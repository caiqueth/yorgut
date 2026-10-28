import React from "react";
import { useMe } from "api/me/queries";
import { HomeMenu } from "components/HomeMenu";
import { PageWithSidebar } from "components/PageWithSidebar";
import { Scrapbook } from "components/Scrapbook";
import { UserProfilePicture } from "components/UserProfilePicture";

export const MyScrapbook: React.FC = () => {
  const { data: user } = useMe();
  if (!user) return null;

  return (
    <PageWithSidebar picture={<UserProfilePicture userId={user.id} />} menu={<HomeMenu />}>
      <Scrapbook userId={user.id} />
    </PageWithSidebar>
  );
};
