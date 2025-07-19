import { useMeMyCommunities } from "api/me/queries";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";

export const MyCommunities: React.FC = () => {
  const { data, isLoading } = useMeMyCommunities();

  const items =
    data?.communities.map((c) => (
      <PicAndTextHorizontal
        title={c.name}
        text={`${c.membershipCount} membros`}
        url={`/community/${c.id}`}
        imageSrc={c.profilePicture}
        imageSize={64}
      />
    )) ?? [];

  return (
    <BoxWithTitle title="minhas comunidades">
      <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} />
    </BoxWithTitle>
  )
}