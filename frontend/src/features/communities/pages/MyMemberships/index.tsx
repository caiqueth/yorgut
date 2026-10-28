import React from "react";
import { useMeMemberships } from "api/me/queries";
import { CommunityRow, CommunityRowSkeleton } from "components/CommunityRow";
import { PagedList } from "components/PagedList";
import { LIST_PAGE_SIZE } from "helpers/paging";
import { usePage } from "hooks/usePage";

// The communities the logged user is a member of, by name.
export const MyMemberships: React.FC = () => {
  const [page, setPage] = usePage();
  const { data, isLoading, isFetching } = useMeMemberships(page, LIST_PAGE_SIZE);

  return (
    <PagedList
      title={"comunidades que participo"}
      count={data?.count}
      items={data?.communities.map((c) => <CommunityRow key={c.id} community={c} />) ?? []}
      isLoading={isLoading}
      isFetching={isFetching}
      skeleton={<CommunityRowSkeleton />}
      empty={"Parece que você ainda não participa de nenhuma comunidade."}
      page={page}
      onPageChange={setPage}
    />
  );
};
