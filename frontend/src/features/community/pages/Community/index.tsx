import React from "react";
import { Outlet } from "react-router-dom";
import { PageWithSidebar } from "components/PageWithSidebar";
import { CommunityProfilePicture } from "features/community/components/CommunityProfilePicture";
import { ProfileMenu } from "features/community/components/ProfileMenu";
import { useIdParam } from "hooks/useIdParam";

export const Community: React.FC = () => {
  const communityId = useIdParam("communityId");

  return (
    <PageWithSidebar
      picture={<CommunityProfilePicture communityId={communityId} />}
      menu={<ProfileMenu communityId={communityId} />}
    >
      <Outlet />
    </PageWithSidebar>
  );
};
