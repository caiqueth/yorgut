import React from "react";
import { useUser, useUserFriends } from "api/user/queries";
import { InlineSkeleton } from "components/InlineSkeleton";
import { PagedList } from "components/PagedList";
import { UserRow, UserRowSkeleton } from "components/UserRow";
import { LIST_PAGE_SIZE } from "helpers/paging";
import { useIdParam } from "hooks/useIdParam";
import { usePage } from "hooks/usePage";

// A profile's friends, most recently seen first.
export const UserFriendsList: React.FC = () => {
  const userId = useIdParam("userId");
  const [page, setPage] = usePage(userId);
  const { data, isLoading, isFetching } = useUserFriends(userId, page, LIST_PAGE_SIZE);
  const { data: user } = useUser(userId);

  return (
    <PagedList
      title={<>amigos de {user?.firstName ?? <InlineSkeleton width={"6ch"} />}</>}
      count={data?.count}
      items={data?.friends.map((f) => <UserRow key={f.id} user={f} />) ?? []}
      isLoading={isLoading}
      isFetching={isFetching}
      skeleton={<UserRowSkeleton />}
      empty={"Ainda não tem nenhum amigo."}
      page={page}
      onPageChange={setPage}
    />
  );
};
