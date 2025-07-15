import { useMeMemberships } from "api/me/queries";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";

interface IProps { }
export const MyMemberships: React.FC<IProps> = () => {
  const { data } = useMeMemberships();

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
    <BoxWithTitle title={"comunidades que participo"}>
      <StripedGrid items={items} elementsPerRow={1} />
    </BoxWithTitle>
  );
}