import React from "react";
import { useUserMemberships } from "api/user/queries";
import { BoxWithTitle } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { PicAndText } from "components/PicAndText";

interface ICommunitiesListProps {
  userId: number;
}
export const CommunitiesList: React.FC<ICommunitiesListProps> = ({
  userId,
}) => {
  const { data: memberships, isLoading } = useUserMemberships(userId);

  const items =
    memberships?.communities.map((c) => (
      <PicAndText
        imageSrc={c.profilePicture}
        url={`/community/${c.id}`}
        text={`${c.name} (${c.membershipCount ?? 0})`}
      />
    )) ?? [];

  return (
    <BoxWithTitle title={`comunidades (${memberships?.membershipCount})`}>
      <StripedGrid items={items} elementsPerRow={3} isLoading={isLoading} />
    </BoxWithTitle>
  );
};
