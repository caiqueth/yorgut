import React from "react";

import { Box, Stack, Typography } from "@mui/material";
import { useSetFan, useUserFans } from "api/userFan/queries";
import { useUpsertUserReview, useUserReviews } from "api/userReview/queries";
import { IUserReview } from "api/userReview/types";

const MAX_RATING = 3;

const traits: { key: keyof IUserReview; label: string; icon: string }[] = [
  { key: "trustworthy", label: "confiável", icon: "😀" },
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
      {Array.from({ length: MAX_RATING }, (_, idx) => {
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
}
const FanStar: React.FC<IFanStarProps> = ({ userId }) => {
  const { data } = useUserFans(userId);
  const { mutate: setFan, isPending } = useSetFan(userId);

  const isFan = data?.amIFan ?? false;
  // Removing a star is always allowed; adding one needs a friendship.
  const disabled = !data || isPending || (!isFan && !data.canFan);

  return (
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
  );
};

interface IFriendRatingProps {
  userId: number;
}
export const FriendRating: React.FC<IFriendRatingProps> = ({ userId }) => {
  const { data } = useUserReviews(userId);
  const { mutate: upsertReview, isPending } = useUpsertUserReview(userId);

  const myReview: IUserReview = data?.myReview ?? { cool: 0, sexy: 0, trustworthy: 0 };
  const disabled = !data?.canReview || isPending;

  const handleChange = (key: keyof IUserReview, value: number) => {
    upsertReview({ ...myReview, [key]: value });
  };

  return (
    <Stack direction={"row"} spacing={2} alignItems={"center"} flexWrap={"wrap"}>
      <Stack direction={"row"} spacing={0.5} alignItems={"center"}>
        <Typography variant={"body2"}>fã</Typography>
        <FanStar userId={userId} />
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
    </Stack>
  );
};
