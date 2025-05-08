import React from "react";
import { useCommunity } from "api/community/queries";

import { Box, Paper, Stack, Typography } from "@mui/material";

interface ICommunityHeaderProps {
  communityId: number;
}
export const ProfileHeader: React.FC<ICommunityHeaderProps> = ({
  communityId,
}) => {
  const { data: community } = useCommunity(communityId);

  return (
    <Box>
      <Paper sx={{ backgroundColor: "white", p: 1 }}>
        <Stack
          direction={"row"}
          sx={{ width: "100%", padding: 0 }}
          alignItems={"end"}
        >
          <Typography
            variant={"h5"}
            sx={{ "&::first-letter": { fontSize: "200%" } }}
          >
            {community?.name ?? "Não encontrada"}
          </Typography>
        </Stack>
      </Paper>
    </Box>
  );
};
