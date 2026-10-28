import React from "react";
import { PicAndTextHorizontal, PicAndTextHorizontalSkeleton } from "components/PicAndTextHorizontal";
import { memberCountText } from "helpers/counts";

interface ICommunityRowProps {
  community: {
    id: number;
    name: string;
    membershipCount: number;
    profilePicture?: string | null;
  };
}

// A community in a list: picture, name and member count, linking to it.
export const CommunityRow: React.FC<ICommunityRowProps> = ({ community }) => (
  <PicAndTextHorizontal
    title={community.name}
    text={memberCountText(community.membershipCount)}
    url={`/community/${community.id}`}
    imageSrc={community.profilePicture}
    imageSize={64}
  />
);

export const CommunityRowSkeleton: React.FC = () => (
  <PicAndTextHorizontalSkeleton imageSize={64} lineWidths={[160, 80]} />
);
