import { Box, Stack, Typography } from "@mui/material";
import { useMeAmMember } from "api/me/queries";
import React from "react";
import { JoinButton } from "features/community/components/JoinButton";

interface IProps {
  communityId: number;
  // What a non-member can't do, e.g. "criar tópicos".
  action: string;
  children?: React.ReactNode;
}
// Shows children only to members of the community. The backend enforces the
// same rule; this just explains it and offers a way to join (or to ask, in a
// private community).
export const MembersOnly: React.FC<IProps> = ({ communityId, action, children }) => {
  const { data, isLoading } = useMeAmMember(communityId);

  if (isLoading) return null;
  if (data?.amIMember) return <>{children}</>;

  return (
    <Box sx={{ p: 1 }}>
      <Stack direction={"row"} alignItems={"center"} gap={1} flexWrap={"wrap"}>
        <Typography variant={"body2"}>
          🔒 Só membros da comunidade podem {action}.
        </Typography>
        <JoinButton communityId={communityId} />
      </Stack>
    </Box>
  );
};
