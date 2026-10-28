import React from "react";
import { useCommunityMembers } from "api/community/queries";
import { PagedList } from "components/PagedList";
import { UserRow, UserRowSkeleton } from "components/UserRow";
import { LIST_PAGE_SIZE } from "helpers/paging";
import { useIdParam } from "hooks/useIdParam";
import { usePage } from "hooks/usePage";

// A community's members, most recently seen first.
export const CommunityMembers: React.FC = () => {
  const communityId = useIdParam("communityId");
  const [page, setPage] = usePage(communityId);
  const { data, isLoading, isFetching } = useCommunityMembers(communityId, page, LIST_PAGE_SIZE);

  return (
    <PagedList
      title={"membros da comunidade"}
      count={data?.count}
      items={data?.members.map((m) => <UserRow key={m.id} user={m} />) ?? []}
      isLoading={isLoading}
      isFetching={isFetching}
      skeleton={<UserRowSkeleton />}
      page={page}
      onPageChange={setPage}
    />
  );
};
