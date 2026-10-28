import React from "react";
import { useUserFanList } from "api/userFan/queries";
import { useUser } from "api/user/queries";
import { InlineSkeleton } from "components/InlineSkeleton";
import { PagedList } from "components/PagedList";
import { UserRow, UserRowSkeleton } from "components/UserRow";
import { LIST_PAGE_SIZE } from "helpers/paging";
import { useIdParam } from "hooks/useIdParam";
import { usePage } from "hooks/usePage";

// The users who are fans of a profile, the newest first.
export const UserFansList: React.FC = () => {
  const userId = useIdParam("userId");
  const [page, setPage] = usePage(userId);
  const { data, isLoading, isFetching } = useUserFanList(userId, page, LIST_PAGE_SIZE);
  const { data: user } = useUser(userId);

  return (
    <PagedList
      title={<>fãs de {user?.firstName ?? <InlineSkeleton width={"6ch"} />}</>}
      count={data?.count}
      items={data?.fans.map((f) => <UserRow key={f.id} user={f} />) ?? []}
      isLoading={isLoading}
      isFetching={isFetching}
      skeleton={<UserRowSkeleton />}
      empty={"Ainda não tem nenhum fã."}
      page={page}
      onPageChange={setPage}
    />
  );
};
