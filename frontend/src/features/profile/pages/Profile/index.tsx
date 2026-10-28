import React from "react";
import { Outlet } from "react-router-dom";
import { PageWithSidebar } from "components/PageWithSidebar";
import { UserProfilePicture } from "components/UserProfilePicture";
import { ProfilePreviewProvider } from "contexts/ProfilePreviewContext";
import { PreviewBanner } from "features/profile/components/PreviewBanner";
import { ProfileMenu } from "features/profile/components/ProfileMenu";
import { useIdParam } from "hooks/useIdParam";

export const Profile: React.FC = () => {
  const userId = useIdParam("userId");

  // Keyed by the profile, so opening another profile ends a preview.
  return (
    <ProfilePreviewProvider key={userId}>
      <PageWithSidebar picture={<UserProfilePicture userId={userId} />} menu={<ProfileMenu userId={userId} />}>
        <PreviewBanner />
        <Outlet />
      </PageWithSidebar>
    </ProfilePreviewProvider>
  );
};
