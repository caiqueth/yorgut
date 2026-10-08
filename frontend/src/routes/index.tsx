import { Home } from "pages/Home";
import { Profile } from "pages/Profile";

import { ProtectedRoute } from "./ProtectedRoute";
import { Community } from "pages/Community";
import { NotFound } from "pages/NotFound";
import { BadBadServer } from "pages/BadBadServer";
import { CommunityMembers } from "pages/Community/Members";
import { EditCommunity } from "pages/Community/Edit";
import { UserFriendsList } from "pages/Profile/Friends";
import { UserCommunitiesList } from "pages/Profile/Communities";
import { UserPhotoAlbumsList } from "pages/Profile/PhotoAlbums";
import { UserScrapbook } from "pages/Profile/Scrapbook";
import { UserVideosList } from "pages/Profile/Videos";
import { UserTestimonialsList } from "pages/Profile/Testimonials";
import { EditProfile } from "pages/Profile/Edit";
import { Friends } from "pages/Friends";
import { Communities } from "pages/Communities";
import { Search } from "pages/Search";
import { MyAccount } from "pages/MyAccount";
import { Help } from "pages/Help";
import { WhatsNew } from "pages/WhatsNew";
import { Media } from "pages/Media";
import { MyScrapbook } from "pages/MyScrapbook";
import { Login } from "pages/Login";
import { MainPage as UserProfileMainPage } from "pages/Profile/MainPage";
import { MainPage as CommunityMainPage } from "pages/Community/MainPage";
import { Forum } from "pages/Community/Forum";
import { Topic } from "pages/Community/Topic";
import { NewTopic } from "pages/Community/NewTopic";
import { NewCommunity } from "pages/Communities/NewCommunity";
import { MyMemberships } from "pages/Communities/MyMemberships";
import { MyCommunities } from "pages/Communities/MyCommunities";
import { DiscoverCommunities } from "pages/Communities/Discover";

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
    children: [
      {
        index: true,
        element: (
          <ProtectedRoute>
            <MyMemberships />
          </ProtectedRoute>
        )
      },
      {
        path: "new",
        element: (
          <ProtectedRoute>
            <NewCommunity />
          </ProtectedRoute>
        )
      },
      {
        path: "owned",
        element: (
          <ProtectedRoute>
            <MyCommunities />
          </ProtectedRoute>
        )
      },
      {
        path: "discover",
        element: (
          <ProtectedRoute>
            <DiscoverCommunities />
          </ProtectedRoute>
        )
      }
    ]
  },
  {
    path: "/scraps",
    element: (
      <ProtectedRoute>
        <MyScrapbook />
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
    path: "/my-account",
    element: (
      <ProtectedRoute>
        <MyAccount />
      </ProtectedRoute>
    ),
  },
  {
    path: "/help",
    element: (
      <ProtectedRoute>
        <Help />
      </ProtectedRoute>
    ),
  },
  {
    path: "/whats-new",
    element: (
      <ProtectedRoute>
        <WhatsNew />
      </ProtectedRoute>
    ),
  },
  {
    path: "/media",
    element: (
      <ProtectedRoute>
        <Media />
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
        element: (
          <ProtectedRoute>
            <UserProfileMainPage />
          </ProtectedRoute>
        )
      },
      {
        path: "scraps",
        element: (
          <ProtectedRoute>
            <UserScrapbook />
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
        path: "edit",
        element: (
          <ProtectedRoute>
            <EditProfile />
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
      { path: "not-allowed", element: <BadBadServer /> },
      { path: "disconnected", element: <BadBadServer /> },
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
            <CommunityMainPage />
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
        path: "edit",
        element: (
          <ProtectedRoute>
            <EditCommunity />
          </ProtectedRoute>
        ),
      },
      {
        path: "forum",
        element: (
          <ProtectedRoute>
            <Forum />
          </ProtectedRoute>
        ),
      },
      {
        path: "forum/topic/:topicId",
        element: (
          <ProtectedRoute>
            <Topic />
          </ProtectedRoute>
        ),
      },
      {
        path: "forum/new-topic",
        element: (
          <ProtectedRoute>
            <NewTopic />
          </ProtectedRoute>
        ),
      },
    ],
  },
  {
    path: "/login",
    element: <Login />,
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
