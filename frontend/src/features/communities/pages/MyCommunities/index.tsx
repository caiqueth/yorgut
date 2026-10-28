import React from "react";
import { useMeMyCommunities } from "api/me/queries";
import { CommunityRow, CommunityRowSkeleton } from "components/CommunityRow";
import { PagedList } from "components/PagedList";
import { LIST_PAGE_SIZE } from "helpers/paging";
import { usePage } from "hooks/usePage";

// The communities the logged user created, by name.
export const MyCommunities: React.FC = () => {
  const [page, setPage] = usePage();
  const { data, isLoading, isFetching } = useMeMyCommunities(page, LIST_PAGE_SIZE);

  return (
    <PagedList
      title={"minhas comunidades"}
      count={data?.count}
      items={data?.communities.map((c) => <CommunityRow key={c.id} community={c} />) ?? []}
      isLoading={isLoading}
      isFetching={isFetching}
      skeleton={<CommunityRowSkeleton />}
      empty={"Parece que você ainda não possui nenhuma comunidade."}
      page={page}
      onPageChange={setPage}
    />
  );
};
