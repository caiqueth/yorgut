import React from "react";
import { BoxWithTitle } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { PicAndText } from "components/PicAndText";
import { useMeMemberships } from "api/me/queries";
import { useUserMemberships } from "api/user/queries";
import { FindUserMembershipsResponse } from "api/generated/model";

interface ICommunitiesListProps {
  // Whose communities to show; omit for the logged user ("minhas comunidades").
  userId?: number;
}

// The communities box of the home page (logged user) and of a profile page.
export const CommunitiesList: React.FC<ICommunitiesListProps> = ({ userId }) =>
  userId === undefined ? <MyCommunitiesList /> : <UserCommunitiesList userId={userId} />;

const MyCommunitiesList: React.FC = () => {
  const { data: memberships, isLoading } = useMeMemberships();
  return (
    <CommunitiesGrid
      title={`minhas comunidades (${memberships?.membershipCount})`}
      memberships={memberships}
      isLoading={isLoading}
    />
  );
};

const UserCommunitiesList: React.FC<{ userId: number }> = ({ userId }) => {
  const { data: memberships, isLoading } = useUserMemberships(userId);
  return (
    <CommunitiesGrid
      title={`comunidades (${memberships?.membershipCount})`}
      memberships={memberships}
      isLoading={isLoading}
    />
  );
};

interface ICommunitiesGridProps {
  title: string;
  memberships?: FindUserMembershipsResponse;
  isLoading: boolean;
}

const CommunitiesGrid: React.FC<ICommunitiesGridProps> = ({ title, memberships, isLoading }) => {
  const items =
    memberships?.communities.map((c) => (
      <PicAndText
        imageSrc={c.profilePicture}
        url={`/community/${c.id}`}
        text={`${c.name} (${c.membershipCount ?? 0})`}
      />
    )) ?? [];

  return (
    <BoxWithTitle title={title}>
      <StripedGrid items={items} elementsPerRow={3} isLoading={isLoading} />
    </BoxWithTitle>
  );
};
