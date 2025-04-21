import { Home } from "pages/Home";
import { Profile } from "pages/Profile";

import { ProtectedRoute } from "./ProtectedRoute";
import { Community } from "pages/Community";
import { NotFound } from "pages/NotFound";
import { BadBadServer } from "pages/BadBadServer";
import { CommunityMembers } from "pages/Community/Members";
import { ProfileDetails as CommunityProfileDetails } from "pages/Community/ProfileDetails";
import { ProfileDetails as UserProfileDetails } from "pages/Profile/ProfileDetails";
import { UserFriendsList } from "pages/Profile/Friends";
import { UserCommunitiesList } from "pages/Profile/Communities";
import { UserPhotoAlbumsList } from "pages/Profile/PhotoAlbums";

export const Routes = [
  {
    path: "/",
    element: (
      <ProtectedRoute>
        <Home />
      </ProtectedRoute>
    ),
  },
  {
    path: "/profile/:userId",
    element: (
      <ProtectedRoute>
        <Profile />
      </ProtectedRoute>
    ),
    children: [
      {
        index: true,
        element: <UserProfileDetails />,
      },
      {
        path: "friends",
        element: (
          <ProtectedRoute>
            <UserFriendsList />
          </ProtectedRoute>
        ),
      },
      {
        path: "communities",
        element: (
          <ProtectedRoute>
            <UserCommunitiesList />
          </ProtectedRoute>
        ),
      },
      {
        path: "albums",
        element: (
          <ProtectedRoute>
            <UserPhotoAlbumsList />
          </ProtectedRoute>
        ),
        children: [
          {
            path: ":albumId",
            element: (
              <ProtectedRoute>
                <></>
              </ProtectedRoute>
            ),
            children: [
              {
                path: ":photoId",
                element: (
                  <ProtectedRoute>
                    <></>
                  </ProtectedRoute>
                ),
              },
            ],
          },
        ],
      },
      { path: "not-allowed", element: <></> },
      { path: "disconnected", element: <></> },
    ],
  },
  {
    path: "/community/:communityId",
    element: (
      <ProtectedRoute>
        <Community />
      </ProtectedRoute>
    ),
    children: [
      {
        index: true,
        element: (
          <ProtectedRoute>
            <CommunityProfileDetails />
          </ProtectedRoute>
        ),
      },
      {
        path: "members",
        element: (
          <ProtectedRoute>
            <CommunityMembers />
          </ProtectedRoute>
        ),
      },
      {
        path: "topics",
        element: (
          <ProtectedRoute>
            <></>
          </ProtectedRoute>
        ),
        children: [
          {
            path: ":topicId",
            element: (
              <ProtectedRoute>
                <></>
              </ProtectedRoute>
            ),
          },
        ],
      },
    ],
  },
  {
    path: "/error",
    element: <NotFound />,
  },
  {
    path: "/not-found",
    element: <NotFound />,
  },
  {
    path: "/bad-server",
    element: <BadBadServer />,
  },
  {
    path: "*",
    element: <NotFound />,
  },
];
