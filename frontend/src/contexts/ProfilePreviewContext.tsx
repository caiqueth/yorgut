import React, { createContext, useContext, useState } from "react";
import { useLocation } from "react-router-dom";

// "Ver como visitante": the user looks at their own profile as people who
// aren't their friends see it. While it's on, the profile's requests ask the
// API for that view (?as=stranger), so the preview is exactly what the server
// sends a non-friend, and owner-only controls (edit, delete, upload) hide.
interface IProfilePreview {
  previewing: boolean;
  setPreviewing: (on: boolean) => void;
}

const ProfilePreviewContext = createContext<IProfilePreview>({
  previewing: false,
  setPreviewing: () => {},
});

// Wraps a profile page. Start the preview from elsewhere by navigating to the
// profile with state { preview: true }; a new profile (key) starts without it.
export const ProfilePreviewProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const location = useLocation();
  const [previewing, setPreviewing] = useState<boolean>(!!location.state?.preview);

  return (
    <ProfilePreviewContext.Provider value={{ previewing, setPreviewing }}>
      {children}
    </ProfilePreviewContext.Provider>
  );
};

// Outside a profile page (home, scrapbook...) it is always off.
export const useProfilePreview = () => useContext(ProfilePreviewContext);

// The query param of profile reads: ?as=stranger while previewing.
export const usePreviewParams = () => {
  const { previewing } = useProfilePreview();
  return previewing ? { as: "stranger" } : {};
};
