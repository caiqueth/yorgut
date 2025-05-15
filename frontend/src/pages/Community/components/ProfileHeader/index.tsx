import React from "react";
import { useCommunity } from "api/community/queries";
import communityTextDecoration from "assets/images/community_text_decoration.svg"
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
          <Box sx={{
            display: "flex",
          }}>
            <Box sx={{
              width: 45,
              height: 45,
              backgroundImage: `url(${communityTextDecoration})`,
              backgroundSize: "cover",
              backgroundRepeat: "no-repeat",
            }}>
            </Box>
            <Typography
              variant={"h5"}
              sx={{
                "&::first-letter": {
                  fontSize: 72,
                  fontFamily: "Gymnastics",
                },
                lineHeight: 1,
                ml: -3
              }}
            >
              {community?.name ?? "Não encontrada"}
            </Typography>
          </Box>
        </Stack>
      </Paper>
    </Box >
  );
};
