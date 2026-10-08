import { Box, Stack, Typography } from "@mui/material";
import { useJoinCommunity, useMeAmMember } from "api/me/queries";
import { Button } from "components/atoms/Button";
import React from "react";

interface IProps {
  communityId: number;
  // What a non-member can't do, e.g. "criar tópicos".
  action: string;
  children: React.ReactNode;
}
// Shows children only to members of the community. The backend enforces the
// same rule; this just explains it and offers a way to join.
export const MembersOnly: React.FC<IProps> = ({ communityId, action, children }) => {
  const { data, isLoading } = useMeAmMember(communityId);
  const { mutate: joinCommunity, isPending } = useJoinCommunity();

  if (isLoading) return null;
  if (data?.amIMember) return <>{children}</>;

  return (
    <Box sx={{ p: 1 }}>
      <Stack direction={"row"} alignItems={"center"} gap={1} flexWrap={"wrap"}>
        <Typography variant={"body2"}>
          🔒 Só membros da comunidade podem {action}.
        </Typography>
        <Button onClick={() => joinCommunity({ communityId })} loading={isPending}>
          participar
        </Button>
      </Stack>
    </Box>
  );
};
