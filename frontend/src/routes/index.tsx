import { Home } from "features/home/pages/Home";
import { Profile } from "features/profile/pages/Profile";

import { ProtectedRoute } from "routes/ProtectedRoute";
import { Community } from "features/community/pages/Community";
import { NotFound } from "features/errors/pages/NotFound";
import { BadBadServer } from "features/errors/pages/BadBadServer";
import { CommunityMembers } from "features/community/pages/Members";
import { EditCommunity } from "features/community/pages/Edit";
import { UserFriendsList } from "features/profile/pages/Friends";
import { UserCommunitiesList } from "features/profile/pages/Communities";
import { UserPhotoAlbumsList } from "features/profile/pages/PhotoAlbums";
import { UserScrapbook } from "features/profile/pages/Scrapbook";
import { UserVideosList } from "features/profile/pages/Videos";
import { UserTestimonialsList } from "features/profile/pages/Testimonials";
import { EditProfile } from "features/profile/pages/Edit";
import { Friends } from "features/friends/pages/Friends";
import { Communities } from "features/communities/pages/Communities";
import { Search } from "features/search/pages/Search";
import { MyAccount } from "features/account/pages/MyAccount";
import { Help } from "features/info/pages/Help";
import { WhatsNew } from "features/info/pages/WhatsNew";
import { Media } from "features/info/pages/Media";
import { MyScrapbook } from "features/scrapbook/pages/MyScrapbook";
import { Login } from "features/auth/pages/Login";
import { MainPage as UserProfileMainPage } from "features/profile/pages/MainPage";
import { MainPage as CommunityMainPage } from "features/community/pages/MainPage";
import { Forum } from "features/community/pages/Forum";
import { Topic } from "features/community/pages/Topic";
import { NewTopic } from "features/community/pages/NewTopic";
import { NewCommunity } from "features/communities/pages/NewCommunity";
import { MyMemberships } from "features/communities/pages/MyMemberships";
import { MyCommunities } from "features/communities/pages/MyCommunities";
import { DiscoverCommunities } from "features/communities/pages/Discover";

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
