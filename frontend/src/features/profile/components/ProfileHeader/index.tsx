import React from "react";
import { useUser, useUserProfileSection } from "api/user/queries";
import { useMe } from "api/me/queries";
import { Button } from "components/Button";
import { useProfilePreview } from "contexts/ProfilePreviewContext";
import profileTextDecoration from "assets/images/profile_text_decoration.svg"
import { Link } from "react-router-dom";

import { Box, Paper, Skeleton, Stack, Typography } from "@mui/material";
import { Repeater } from "components/Repeater";
import { useUserReviews } from "api/userReview/queries";
import { useUserFans } from "api/userFan/queries";
import { usePhotos } from "api/photo/queries";
import { useScraps } from "api/scrap/queries";
import { SCRAP_PAGE_SIZE } from "helpers/paging";
import { ratingIcons, ratingPercent } from "helpers/ratings";
import { InlineSkeleton, loadingTextSx, TextPlaceholder } from "components/InlineSkeleton";

interface IStatProps {
  icon: string;
  // Undefined while loading.
  value?: number;
  label: string;
  // Where the count leads, e.g. the scrapbook; omit for plain text.
  to?: string;
}
// One count under the name ("📝 12 recados"); only the number is a skeleton
// until it loads.
const Stat: React.FC<IStatProps> = ({ icon, value, label, to }) => {
  const text = <>{value ?? <InlineSkeleton width={"2ch"} />} {label}</>;
  return (
    <Typography variant={"body2"}>
      {icon} {to ? <Link to={to}>{text}</Link> : text}
    </Typography>
  );
};

interface IProfileHeaderProps {
  userId: number;
  // The box below continues this one (ProfileDetails' joinedAbove): no
  // bottom border, corners or shadow.
  joinedBelow?: boolean;
}
export const ProfileHeader: React.FC<IProfileHeaderProps> = ({ userId, joinedBelow }) => {
  const { data: user } = useUser(userId);
  const { data: reviews, isLoading: isLoadingReviews } = useUserReviews(userId);
  const { data: fans } = useUserFans(userId);
  const { data: me } = useMe();
  const { previewing, setPreviewing } = useProfilePreview();
  // Counts of sections kept to friends aren't shown to others.
  const album = useUserProfileSection(userId, "album");
  const scrapbook = useUserProfileSection(userId, "scraps");
  const { data: photos } = usePhotos(userId, album.canSee);
  // Same query as the scrapbook's first page, so both share one request.
  const { data: scraps } = useScraps(userId, 1, SCRAP_PAGE_SIZE, scrapbook.canSee);

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
        <Stack direction={"row"} sx={{ width: "100%", padding: 0 }} alignItems={"end"}>
          {/* The circle sits at the top, behind the big first letter (which
              ml: -3 pulls over it); aligned to the bottom it drops below the
              letter's top. */}
          <Box sx={{
            display: "flex",
            alignItems: "start",
          }}>
            <Box sx={{
              width: 45,
              height: 45,
              flexShrink: 0,
              backgroundImage: `url(${profileTextDecoration})`,
              backgroundSize: "cover",
              backgroundRepeat: "no-repeat",
            }}>
            </Box>
            <Typography
              variant={"h5"}
              sx={{
                "&::first-letter": {
                  fontSize: 72,
                  ml: -3
                },
                lineHeight: 1,
                ...(user ? {} : loadingTextSx),
              }}>
              {user ? `${user.firstName} ${user.lastName}` : <TextPlaceholder />}
            </Typography>
          </Box>
        </Stack>
        <Stack direction={"row"} spacing={1}>
          {(scrapbook.isLoading || scrapbook.canSee) && (
            <Stat
              icon={"📝"}
              value={scraps?.count}
              label={"recados"}
              to={`/profile/${userId}/scraps`}
            />
          )}
          {(album.isLoading || album.canSee) && (
            <Stat
              icon={"📷"}
              value={photos?.count}
              label={"fotos"}
              to={`/profile/${userId}/albums`}
            />
          )}
          <Stat
            icon={"⭐"}
            value={fans?.fanCount}
            label={"fãs"}
            to={`/profile/${userId}/fans`}
          />
        </Stack>
        <Stack direction={"row"} justifyContent={"space-between"} alignItems={"center"} flexWrap={"wrap"} gap={1}>
          {/* Keeps its height while empty, so the page doesn't move when ratings arrive. */}
          <Stack direction={"row"} spacing={0} sx={{ minHeight: 24 }} alignItems={"center"}>
            {isLoadingReviews
              ? <Skeleton variant={"text"} width={96} />
              : (
                <>
                  {([
                    ["🧊", "legal", reviews?.cool],
                    ["🙂", "confiável", reviews?.trustworthy],
                    ["❤️", "sexy", reviews?.sexy],
                  ] as const).map(([icon, trait, average = 0]) => (
                    <Repeater
                      key={trait}
                      element={icon}
                      times={ratingIcons(average)}
                      title={ratingPercent(average)}
                      label={`${trait}: ${ratingPercent(average)}`}
                    />
                  ))}
                </>
              )}
          </Stack>
          {/* The user's own profile: see it as people who aren't their friends do. */}
          {me?.id === userId && !previewing && (
            <Button size={"small"} onClick={() => setPreviewing(true)}>👁️ ver como visitante</Button>
          )}
        </Stack>
      </Paper>
    </Box>
  );
};
