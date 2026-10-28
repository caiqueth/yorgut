import {
  useFindUser,
  useFindUserFriends,
  useFindUserMemberships,
  useFindUserProfile,
  useSearchUsers,
} from "api/generated/endpoints/users/users";
import { usePreviewParams } from "contexts/ProfilePreviewContext";

export const useUser = (userId: number) => useFindUser(userId);

// One page of the user's friends, most recently seen first.
export const useUserFriends = (userId: number, page: number, pageSize: number) =>
  useFindUserFriends(userId, { page, pageSize }, { query: { placeholderData: (previous) => previous } });

// One page of the user's communities, by name, or shuffled for the profile
// box (the API keeps the same shuffle all day, so reloads don't reshuffle).
// enabled false skips the request, e.g. for communities kept to friends.
export const useUserMemberships = (
  userId: number,
  page: number,
  pageSize: number,
  random = false,
  enabled = true,
) =>
  useFindUserMemberships(
    userId,
    { page, pageSize, random, ...usePreviewParams() },
    { query: { enabled, placeholderData: (previous) => previous } },
  );

// The profile as the logged user may see it (or, while previewing, as
// non-friends see it): visibility tells which sections they may open.
export const useUserProfile = (userId: number) =>
  useFindUserProfile(userId, usePreviewParams(), { query: { enabled: !!userId } });

// One page of the users whose name contains name (empty = everyone).
export const useUserSearchUsers = (name: string, page: number, pageSize: number) =>
  useSearchUsers(
    { name, page, pageSize },
    { query: { placeholderData: (previous) => previous } },
  );

// Whether the logged user may see a part of userId's profile (see
// visibility in the profile response). isLoading until the profile arrives:
// the section's own request waits for it.
export const useUserProfileSection = (
  userId: number,
  section: "album" | "communities" | "scraps" | "testimonials",
) => {
  const { data, isLoading } = useUserProfile(userId);
  return { canSee: !!data?.visibility[section], isLoading, firstName: data?.firstName };
};
