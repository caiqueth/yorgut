import React from "react";

import { Box, CircularProgress, Stack, Typography } from "@mui/material";
import { useSetFan } from "api/userFan/queries";
import { useUpsertUserReview } from "api/userReview/queries";
import { UserReviewItem } from "api/generated/model";
import { RATING_MAX } from "helpers/ratings";

const traits: { key: keyof UserReviewItem; label: string; icon: string }[] = [
  { key: "trustworthy", label: "confiável", icon: "🙂" },
  { key: "cool", label: "legal", icon: "🧊" },
  { key: "sexy", label: "sexy", icon: "❤️" },
];

interface IRatingIconsProps {
  icon: string;
  value: number;
  disabled: boolean;
  onChange: (value: number) => void;
}
const RatingIcons: React.FC<IRatingIconsProps> = ({ icon, value, disabled, onChange }) => {
  const [hovered, setHovered] = React.useState<number | null>(null);
  const shown = hovered ?? value;

  return (
    <Stack direction={"row"} onMouseLeave={() => setHovered(null)}>
      {Array.from({ length: RATING_MAX }, (_, idx) => {
        const rating = idx + 1;
        return (
          <Box
            key={rating}
            component={"span"}
            onMouseEnter={() => setHovered(rating)}
            // Clicking the current rating again clears it
            onClick={() => !disabled && onChange(rating === value ? 0 : rating)}
            sx={{
              cursor: disabled ? "default" : "pointer",
              opacity: rating <= shown ? 1 : 0.25,
              userSelect: "none",
            }}
          >
            {icon}
          </Box>
        );
      })}
    </Stack>
  );
};

interface IFanStarProps {
  userId: number;
  isFan: boolean;
}
const FanStar: React.FC<IFanStarProps> = ({ userId, isFan }) => {
  const { mutate: setFan, isPending } = useSetFan(userId);
  const disabled = isPending;

  return (
    <Box component={"span"} sx={{ position: "relative", display: "inline-flex" }}>
      <Box
        component={"span"}
        title={isFan ? "deixar de ser fã" : "ser fã"}
        onClick={() => !disabled && setFan(!isFan)}
        sx={{
          cursor: disabled ? "default" : "pointer",
          opacity: isFan ? 1 : 0.25,
          userSelect: "none",
        }}
      >
        ⭐
      </Box>
      {/* Over the star while saving; hidden, not removed, so nothing moves. */}
      <CircularProgress
        size={14}
        sx={{
          position: "absolute",
          inset: 0,
          margin: "auto",
          visibility: isPending ? "visible" : "hidden",
        }}
      />
    </Box>
  );
};

interface IFriendRatingProps {
  userId: number;
  // The logged user's rating of this friend and fan star, from the friends list
  // (GET /me/friends), so rows don't each ask the API.
  myReview: UserReviewItem;
  amIFan: boolean;
}
// Rates a friend of the logged user (only friends can be rated or starred, and
// the friends page only lists friends).
export const FriendRating: React.FC<IFriendRatingProps> = ({ userId, myReview, amIFan }) => {
  const { mutate: upsertReview, isPending } = useUpsertUserReview(userId);
  const disabled = isPending;

  const handleChange = (key: keyof UserReviewItem, value: number) => {
    upsertReview({ ...myReview, [key]: value });
  };

  return (
    <Stack direction={"row"} spacing={2} alignItems={"center"} flexWrap={"wrap"}>
      <Stack direction={"row"} spacing={0.5} alignItems={"center"}>
        <Typography variant={"body2"}>fã</Typography>
        <FanStar userId={userId} isFan={amIFan} />
      </Stack>
      {traits.map(({ key, label, icon }) => (
        <Stack key={key} direction={"row"} spacing={0.5} alignItems={"center"}>
          <Typography variant={"body2"}>{label}</Typography>
          <RatingIcons
            icon={icon}
            value={myReview[key]}
            disabled={disabled}
            onChange={(value) => handleChange(key, value)}
          />
        </Stack>
      ))}
      {/* Hidden, not removed, so the line doesn't change width while saving. */}
      <CircularProgress size={14} sx={{ visibility: isPending ? "visible" : "hidden" }} />
    </Stack>
  );
};
