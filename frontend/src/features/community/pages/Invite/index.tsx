import React from "react";
import { Stack, Typography } from "@mui/material";
import { useInvitableFriends, useInviteToCommunity } from "api/community/queries";
import { useMeAmMember } from "api/me/queries";
import { InvitableFriendItem } from "api/generated/model";
import { Button } from "components/Button";
import { PagedList } from "components/PagedList";
import { UserRow, UserRowSkeleton } from "components/UserRow";
import { MembersOnly } from "features/community/components/MembersOnly";
import { LIST_PAGE_SIZE } from "helpers/paging";
import { useIdParam } from "hooks/useIdParam";
import { usePage } from "hooks/usePage";

// A member's friends, to invite to the community.
export const InviteFriends: React.FC = () => {
  const communityId = useIdParam("communityId");
  const [page, setPage] = usePage(communityId);
  const { data: membership } = useMeAmMember(communityId);
  // Only members may invite; the backend answers 403 to anyone else.
  const { data, isLoading, isFetching } = useInvitableFriends(communityId, page, LIST_PAGE_SIZE, !!membership?.amIMember);
  const { mutate: invite, isPending, variables: inviting } = useInviteToCommunity(communityId);

  const status = (f: InvitableFriendItem) => {
    if (f.member) return <Typography variant={"body2"}>já participa</Typography>;
    if (f.invited) return <Typography variant={"body2"}>convidado(a) ✉️</Typography>;
    return (
      <Button onClick={() => invite(f.id)} loading={isPending && inviting === f.id}>
        convidar
      </Button>
    );
  };

  const items = (data?.friends ?? []).map((f) => (
    <Stack key={f.id} direction={"row"} justifyContent={"space-between"} alignItems={"center"} gap={1} flexWrap={"wrap"}>
      <UserRow user={f} />
      {status(f)}
    </Stack>
  ));

  return (
    <MembersOnly communityId={communityId} action={"convidar amigos"}>
      <PagedList
        title={"convidar amigos"}
        count={data?.count}
        items={items}
        isLoading={isLoading}
        isFetching={isFetching}
        skeleton={<UserRowSkeleton />}
        empty={"Você ainda não tem amigos para convidar."}
        page={page}
        onPageChange={setPage}
      />
    </MembersOnly>
  );
};
