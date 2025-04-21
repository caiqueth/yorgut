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
import { UserScrapList } from "pages/Profile/Scraps";
import { UserVideosList } from "pages/Profile/Videos";
import { UserTestimonialsList } from "pages/Profile/Testimonials";
import { RelatedCommunities } from "pages/Community/Related";
import { Friends } from "pages/Friends";
import { Communities } from "pages/Communities";
import { Search } from "pages/Search";

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
    path: "/friends",
    element: (
      <ProtectedRoute>
        <Friends />
      </ProtectedRoute>
    ),
  },
  {
    path: "/communities",
    element: (
      <ProtectedRoute>
        <Communities />
      </ProtectedRoute>
    ),
  },
  {
    path: "/search",
    element: (
      <ProtectedRoute>
        <Search />
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
        path: "scraps",
        element: (
          <ProtectedRoute>
            <UserScrapList />
          </ProtectedRoute>
        ),
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
        path: "videos",
        element: (
          <ProtectedRoute>
            <UserVideosList />
          </ProtectedRoute>
        ),
      },
      {
        path: "testimonials",
        element: (
          <ProtectedRoute>
            <UserTestimonialsList />
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
        path: "related",
        element: (
          <ProtectedRoute>
            <RelatedCommunities />
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
