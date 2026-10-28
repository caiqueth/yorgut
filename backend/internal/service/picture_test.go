package service

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"image"
	"image/png"
	"testing"
	"yorgut/internal/storage"
)

func testPNG(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// TestUploadsStoreKeys checks that uploads go to the object store and the
// database keeps only the key (the response mappers turn it into a URL).
func TestUploadsStoreKeys(t *testing.T) {
	ctx := context.Background()
	repos := newTestRepositories(t)
	store := storage.NewMemory()
	s := NewServices(repos, store)
	owner := createTestUser(t, repos, "owner@test.com")

	t.Run("album photo", func(t *testing.T) {
		if err := s.Photo.UploadPhoto(ctx, owner.ID, testPNG(t), "oi"); err != nil {
			t.Fatal(err)
		}
		photos, err := s.Photo.FindPhotosOfUser(owner.ID, owner.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(photos.Photos) != 1 || !store.Has(photos.Photos[0].URL) {
			t.Fatalf("photos = %+v, want one photo whose URL is a stored key", photos.Photos)
		}
	})

	t.Run("community picture", func(t *testing.T) {
		id, err := s.Community.CreateCommunity(ctx, owner.ID, "Comunidade", "other", "", false, testPNG(t))
		if err != nil {
			t.Fatal(err)
		}
		profile, err := repos.CommunityProfile.FindOfCommunity(id)
		if err != nil {
			t.Fatal(err)
		}
		if profile.ProfilePicture == nil || !store.Has(*profile.ProfilePicture) {
			t.Fatalf("community picture = %v, want a stored key", profile.ProfilePicture)
		}
	})

	t.Run("invalid picture stores nothing", func(t *testing.T) {
		before := store.Len()
		err := s.UserProfile.UpdateProfilePicture(ctx, owner.ID, []byte("not an image"))
		if !errors.Is(err, ErrInvalidPicture) {
			t.Fatalf("err = %v, want ErrInvalidPicture", err)
		}
		if store.Len() != before {
			t.Fatalf("store has %d objects, want %d", store.Len(), before)
		}
	})
}

func TestAlbumLimit(t *testing.T) {
	ctx := context.Background()
	repos := newTestRepositories(t)
	store := storage.NewMemory()
	s := NewServices(repos, store)
	owner := createTestUser(t, repos, "owner@test.com")

	for i := 0; i < MaxPhotosPerAlbum; i++ {
		if err := s.Photo.UploadPhoto(ctx, owner.ID, testPNG(t), ""); err != nil {
			t.Fatalf("photo %d: %v", i+1, err)
		}
	}

	stored := store.Len()
	if err := s.Photo.UploadPhoto(ctx, owner.ID, testPNG(t), ""); !errors.Is(err, ErrAlbumFull) {
		t.Fatalf("photo %d: got %v, want ErrAlbumFull", MaxPhotosPerAlbum+1, err)
	}
	if store.Len() != stored {
		t.Error("the refused photo was stored anyway")
	}

	photos, err := s.Photo.FindPhotosOfUser(owner.ID, owner.ID)
	if err != nil || photos.Count != MaxPhotosPerAlbum || photos.MaxPhotos != 12 {
		t.Fatalf("got count %d, max %d, err %v; want %d and 12", photos.Count, photos.MaxPhotos, err, MaxPhotosPerAlbum)
	}

	// Deleting one makes room again.
	if err := s.Photo.DeletePhoto(ctx, owner.ID, photos.Photos[0].ID); err != nil {
		t.Fatal(err)
	}
	if err := s.Photo.UploadPhoto(ctx, owner.ID, testPNG(t), ""); err != nil {
		t.Errorf("after deleting one: %v", err)
	}
}

// TestReplacedPicturesAreDeleted checks that a picture's file goes away once
// nothing points to it: replaced profile and community pictures, deleted
// album photos.
func TestReplacedPicturesAreDeleted(t *testing.T) {
	ctx := context.Background()
	repos := newTestRepositories(t)
	store := storage.NewMemory()
	s := NewServices(repos, store)
	owner := createTestUser(t, repos, "owner@test.com")

	t.Run("profile picture", func(t *testing.T) {
		// A Google picture is a URL, not a key: replacing it deletes nothing.
		if _, err := repos.UserProfile.UpdatePictureOfUser(owner.ID, "https://lh3.googleusercontent.com/a/x"); err != nil {
			t.Fatal(err)
		}
		pictureKey := func() string {
			t.Helper()
			user, err := repos.User.FindUserWithProfile(owner.ID)
			if err != nil {
				t.Fatal(err)
			}
			return *user.Profile.ProfilePicture
		}
		if err := s.UserProfile.UpdateProfilePicture(ctx, owner.ID, testPNG(t)); err != nil {
			t.Fatal(err)
		}
		first := pictureKey()
		if !store.Has(first) {
			t.Fatal("the new picture is not stored")
		}
		if err := s.UserProfile.UpdateProfilePicture(ctx, owner.ID, testPNG(t)); err != nil {
			t.Fatal(err)
		}
		if store.Has(first) {
			t.Error("the replaced profile picture is still stored")
		}
		if !store.Has(pictureKey()) {
			t.Error("the new picture is not stored")
		}
	})

	t.Run("community picture", func(t *testing.T) {
		id, err := s.Community.CreateCommunity(ctx, owner.ID, "Comunidade", "other", "", false, testPNG(t))
		if err != nil {
			t.Fatal(err)
		}
		profile, err := repos.CommunityProfile.FindOfCommunity(id)
		if err != nil {
			t.Fatal(err)
		}
		old := *profile.ProfilePicture
		if err := s.Community.UpdateCommunityPicture(ctx, owner.ID, id, testPNG(t)); err != nil {
			t.Fatal(err)
		}
		if store.Has(old) {
			t.Error("the replaced community picture is still stored")
		}
		profile, err = repos.CommunityProfile.FindOfCommunity(id)
		if err != nil {
			t.Fatal(err)
		}
		if !store.Has(*profile.ProfilePicture) {
			t.Error("the new community picture is not stored")
		}
	})

	t.Run("album photo", func(t *testing.T) {
		if err := s.Photo.UploadPhoto(ctx, owner.ID, testPNG(t), ""); err != nil {
			t.Fatal(err)
		}
		photos, err := s.Photo.FindPhotosOfUser(owner.ID, owner.ID)
		if err != nil {
			t.Fatal(err)
		}
		key := photos.Photos[0].URL
		if err := s.Photo.DeletePhoto(ctx, owner.ID, photos.Photos[0].ID); err != nil {
			t.Fatal(err)
		}
		if store.Has(key) {
			t.Error("the deleted photo's file is still stored")
		}
	})
}

// pngClaiming returns a 1x1 PNG whose header says it is width x height: the
// limits are checked on the header, before anything is decoded.
func pngClaiming(t *testing.T, width, height uint32) []byte {
	t.Helper()
	data := append([]byte{}, testPNG(t)...)
	// The IHDR chunk follows the 8-byte signature: length(4) type(4) data(13) crc(4).
	ihdr := data[8+8 : 8+8+13]
	binary.BigEndian.PutUint32(ihdr[0:], width)
	binary.BigEndian.PutUint32(ihdr[4:], height)
	binary.BigEndian.PutUint32(data[8+8+13:], crc32.ChecksumIEEE(data[8+4:8+8+13]))
	return data
}

func TestPreparePicture(t *testing.T) {
	t.Run("over 5 MB is refused", func(t *testing.T) {
		_, err := preparePicture(make([]byte, MaxUploadBytes+1))
		if !errors.Is(err, ErrPictureTooLarge) {
			t.Fatalf("err = %v, want ErrPictureTooLarge", err)
		}
	})

	t.Run("over 50 megapixels is refused before decoding", func(t *testing.T) {
		_, err := preparePicture(pngClaiming(t, 10000, 10000))
		if !errors.Is(err, ErrPictureTooBig) {
			t.Fatalf("err = %v, want ErrPictureTooBig", err)
		}
	})

	t.Run("not an image is refused", func(t *testing.T) {
		_, err := preparePicture([]byte("<html>hi</html>"))
		if !errors.Is(err, ErrInvalidPicture) {
			t.Fatalf("err = %v, want ErrInvalidPicture", err)
		}
	})
}
