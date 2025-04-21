import React from "react";

import { useUserProfile } from "api/user/queries";
import { ProfileTable } from "components/organisms/ProfileTable";

interface IProfileDetailsProps {
  userId: number;
}
export const ProfileDetails: React.FC<IProfileDetailsProps> = ({ userId }) => {
  const { data: userProfile } = useUserProfile(userId);

  const data = [
    {
      label: "quem sou eu:",
      value: userProfile?.description,
    },
    {
      label: "relacionamento:",
      value: userProfile?.relationshipStatus,
    },
    {
      label: "interesses:",
      value: userProfile?.interests,
    },
    {
      label: "sobre mim:",
      value: userProfile?.about,
    },
    {
      label: "religião:",
      value: userProfile?.religion,
    },
    {
      label: "humor:",
      value: userProfile?.humor,
    },
    {
      label: "moda:",
      value: userProfile?.fashion,
    },
    {
      label: "mora com:",
      value: userProfile?.livesWith,
    },
    {
      label: "cidade natal:",
      value: userProfile?.birthCity,
    },
    {
      label: "página:",
      value: userProfile?.website,
    },
    {
      label: "paixões:",
      value: userProfile?.passions,
    },
    {
      label: "esportes:",
      value: userProfile?.sports,
    },
    {
      label: "atividades:",
      value: userProfile?.hobbies,
    },
    {
      label: "livros:",
      value: userProfile?.books,
    },
    {
      label: "música:",
      value: userProfile?.music,
    },
    {
      label: "filmes:",
      value: userProfile?.movies,
    },
    {
      label: "culinária:",
      value: userProfile?.cuisine,
    },
  ];

  return <ProfileTable data={data} />;
};
