import {
  useFindUser,
  useFindUserFriends,
  useFindUserMemberships,
  useFindUserProfile,
  useFindUsers,
  useSearchUsers,
} from "api/generated/endpoints/users/users";

export const useUsers = () => useFindUsers();

export const useUser = (userId: number) => useFindUser(userId);

export const useUserFriends = (userId: number) => useFindUserFriends(userId);

export const useUserMemberships = (userId: number) => useFindUserMemberships(userId);

export const useUserProfile = (userId: number) =>
  useFindUserProfile(userId, { query: { enabled: !!userId } });

export const useUserSearchUsers = (name: string, page: number, pageSize: number) =>
  useSearchUsers({ name, page, pageSize });
