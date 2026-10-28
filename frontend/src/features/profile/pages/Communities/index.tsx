import React from "react";
import { useUser, useUserMemberships, useUserProfileSection } from "api/user/queries";
import { BoxWithTitle } from "components/BoxWithTitle";
import { PrivateSection } from "components/PrivateSection";
import { CommunityRow, CommunityRowSkeleton } from "components/CommunityRow";
import { InlineSkeleton } from "components/InlineSkeleton";
import { PagedList } from "components/PagedList";
import { LIST_PAGE_SIZE } from "helpers/paging";
import { useIdParam } from "hooks/useIdParam";
import { usePage } from "hooks/usePage";

// The communities a profile's user is a member of, by name.
export const UserCommunitiesList: React.FC = () => {
  const userId = useIdParam("userId");
  const [page, setPage] = usePage(userId);
  const section = useUserProfileSection(userId, "communities");
  const { data, isLoading, isFetching } = useUserMemberships(userId, page, LIST_PAGE_SIZE, false, section.canSee);
  const { data: user } = useUser(userId);

  if (!section.isLoading && !section.canSee) {
    return (
      <BoxWithTitle title={"comunidades"}>
        <PrivateSection firstName={section.firstName} what={"as comunidades"} />
      </BoxWithTitle>
    );
  }

  return (
    <PagedList
      title={<>comunidades de {user?.firstName ?? <InlineSkeleton width={"6ch"} />}</>}
      count={data?.count}
      items={data?.communities.map((c) => <CommunityRow key={c.id} community={c} />) ?? []}
      // The request waits for the profile, which says whether it may be made.
      isLoading={section.isLoading || isLoading}
      isFetching={isFetching}
      skeleton={<CommunityRowSkeleton />}
      empty={"Ainda não participa de nenhuma comunidade."}
      page={page}
      onPageChange={setPage}
    />
  );
};
