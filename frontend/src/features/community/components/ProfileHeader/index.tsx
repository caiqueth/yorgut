import React from "react";
import { useCommunity } from "api/community/queries";
import communityTextDecoration from "assets/images/community_text_decoration.svg"
import { Box, Paper, Stack, Typography } from "@mui/material";
import { loadingTextSx, TextPlaceholder } from "components/InlineSkeleton";

interface ICommunityHeaderProps {
  communityId: number;
  // The box below continues this one (ProfileDetails' joinedAbove): no
  // bottom border, corners or shadow.
  joinedBelow?: boolean;
}
export const ProfileHeader: React.FC<ICommunityHeaderProps> = ({
  communityId,
  joinedBelow,
}) => {
  const { data: community } = useCommunity(communityId);

  return (
    <Box>
      <Paper sx={{
        backgroundColor: "white",
        p: 1,
        borderRadius: (theme) => theme.yorgut.featuredBoxRadius,
        ...(joinedBelow && {
          borderBottom: 0,
          borderBottomLeftRadius: 0,
          borderBottomRightRadius: 0,
          clipPath: "inset(-4px -4px 0 -4px)",
        }),
      }}>
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
                  fontFamily: "Gymnastics, cursive",
                },
                lineHeight: 1,
                ml: -3,
                ...(community ? {} : loadingTextSx),
              }}
            >
              {community ? community.name : <TextPlaceholder />}
            </Typography>
          </Box>
        </Stack>
      </Paper>
    </Box >
  );
};
