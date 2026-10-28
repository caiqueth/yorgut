import { useIdParam } from "hooks/useIdParam";
import { Box, Dialog, DialogContent, Skeleton, Stack, TextField, Typography } from "@mui/material";
import { useDeletePhoto, usePhotos, useUploadPhoto } from "api/photo/queries";
import { PhotoItem } from "api/generated/model";
import { Button } from "components/Button";
import { BoxWithTitle } from "components/BoxWithTitle";
import { InlineSkeleton } from "components/InlineSkeleton";
import { useMe } from "api/me/queries";
import { useUserProfileSection } from "api/user/queries";
import { PrivateSection } from "components/PrivateSection";
import { useProfilePreview } from "contexts/ProfilePreviewContext";
import { dayjsFromDate } from "helpers/formatting/date";
import { PICTURE_RULES_TEXT } from "helpers/validatePicture";
import { usePicturePicker } from "hooks/usePicturePicker";
import React, { useState } from "react";

// Same limit as the backend (service/photo.go).
const MAX_CAPTION_LENGTH = 200;

export const UserPhotoAlbumsList: React.FC = () => {
  const userId = useIdParam("userId");

  const { data: me } = useMe();
  const { previewing } = useProfilePreview();
  const section = useUserProfileSection(userId, "album");
  const { data, isLoading: isLoadingPhotos, isFetching } = usePhotos(userId, section.canSee);
  // The request waits for the profile, which says whether it may be made.
  const isLoading = section.isLoading || isLoadingPhotos;
  const { mutate: uploadPhoto, isPending: isUploading, error: uploadError, reset: resetUpload } = useUploadPhoto(userId);
  const { mutate: deletePhoto, isPending: isDeleting, variables: deletingId } = useDeletePhoto(userId);

  const { inputRef, picture, preview, error: validationError, pick, clear } = usePicturePicker();
  const [caption, setCaption] = useState("");
  const [opened, setOpened] = useState<PhotoItem | null>(null);

  const count = data?.count ?? 0;
  const maxPhotos = data?.maxPhotos ?? 12;
  const isFull = count >= maxPhotos;
  // Known before the photos load, so the upload box doesn't push the album down.
  const canUpload = data ? data.canUpload : me?.id === userId && !previewing;

  const handlePick = (e: React.ChangeEvent<HTMLInputElement>) => {
    resetUpload();
    pick(e);
  };

  const handleUpload = () => {
    if (!picture) return;
    uploadPhoto({ picture, caption }, {
      onSuccess: () => {
        clear();
        setCaption("");
      },
    });
  };

  const photos = data?.photos ?? [];

  if (!section.isLoading && !section.canSee) {
    return (
      <BoxWithTitle title={"álbum"}>
        <PrivateSection firstName={section.firstName} what={"o álbum"} />
      </BoxWithTitle>
    );
  }

  return (
    <Stack direction={"column"} gap={1}>
      {canUpload && (
        <BoxWithTitle title={"adicionar foto"}>
          <Stack direction={"column"} gap={1} sx={{ p: 1 }}>
            {isFull ? (
              <Typography variant={"body2"}>
                Seu álbum já tem {maxPhotos} fotos. Exclua uma para adicionar outra.
              </Typography>
            ) : (
              <>
                <Typography variant={"body2"}>{PICTURE_RULES_TEXT}</Typography>
                <input ref={inputRef} type={"file"} style={{ maxWidth: "100%" }} accept={"image/png,image/jpeg"} onChange={handlePick} />
                {preview && <img src={preview} height={120} style={{ alignSelf: "start" }} alt={"pré-visualização"} />}
                <TextField
                  label={"legenda (opcional)"}
                  value={caption}
                  onChange={(e) => setCaption(e.target.value)}
                  slotProps={{ htmlInput: { maxLength: MAX_CAPTION_LENGTH } }}
                />
                <Box>
                  <Button onClick={handleUpload} disabled={!picture} loading={isUploading}>
                    enviar foto
                  </Button>
                </Box>
              </>
            )}
            {validationError && (
              <Typography variant={"body2"} color={"error"}>{validationError}</Typography>
            )}
            {uploadError && (
              <Typography variant={"body2"} color={"error"}>
                {uploadError.message}
              </Typography>
            )}
          </Stack>
        </BoxWithTitle>
      )}

      <BoxWithTitle
        title={<>álbum ({data ? `${data.count}/${data.maxPhotos}` : <InlineSkeleton width={"4ch"} />})</>}
        showProgress={isFetching && !isLoading}
      >
        {!isLoading && photos.length === 0 ? (
          <Box sx={{ p: 1 }}>
            <Typography>Nenhuma foto ainda.</Typography>
          </Box>
        ) : (
          <Box
            sx={{
              p: 1,
              display: "grid",
              // Always 3 per row: a full album (12) is 4 rows. minmax(0, …) keeps a
              // long caption from widening its column.
              gridTemplateColumns: "repeat(3, minmax(0, 1fr))",
              gap: 1,
            }}
          >
            {isLoading && Array.from({ length: 3 }, (_, idx) => (
              <Skeleton key={idx} variant={"rectangular"} sx={{ width: "100%", height: "auto", aspectRatio: "1" }} />
            ))}
            {photos.map((photo) => (
              <Stack key={photo.id} direction={"column"} gap={0.5}>
                <Box
                  component={"img"}
                  src={photo.url}
                  alt={photo.caption || "foto"}
                  onClick={() => setOpened(photo)}
                  sx={{
                    width: "100%",
                    aspectRatio: "1",
                    objectFit: "cover",
                    cursor: "pointer",
                    border: "1px solid #ccc",
                  }}
                />
                {photo.caption && (
                  <Typography variant={"body2"} sx={{ overflowWrap: "anywhere" }}>{photo.caption}</Typography>
                )}
                {data?.canUpload && (
                  <Box>
                    <Button
                      onClick={() => deletePhoto(photo.id)}
                      loading={isDeleting && deletingId === photo.id}
                    >
                      excluir
                    </Button>
                  </Box>
                )}
              </Stack>
            ))}
          </Box>
        )}
      </BoxWithTitle>

      <Dialog open={!!opened} onClose={() => setOpened(null)} maxWidth={"lg"}>
        {opened && (
          <DialogContent>
            <Stack direction={"column"} gap={1} alignItems={"center"}>
              <img
                src={opened.url}
                alt={opened.caption || "foto"}
                style={{ maxWidth: "100%", maxHeight: "80vh" }}
              />
              {opened.caption && <Typography>{opened.caption}</Typography>}
              <Typography variant={"body2"}>
                {dayjsFromDate(opened.createdAt, "D [de] MMMM [de] YYYY")}
              </Typography>
            </Stack>
          </DialogContent>
        )}
      </Dialog>
    </Stack>
  );
};
