import { Home } from "pages/Home";
import { Profile } from "pages/Profile";

import { ProtectedRoute } from "./ProtectedRoute";
import { Community } from "pages/Community";

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
        path: "friends",
        element: (
          <ProtectedRoute>
            <></>
          </ProtectedRoute>
        ),
      },
      {
        path: "communities",
        element: (
          <ProtectedRoute>
            <></>
          </ProtectedRoute>
        ),
      },
      {
        path: "albums",
        element: (
          <ProtectedRoute>
            <></>
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
  },
  {
    path: "/community/:communityId",
    element: (
      <ProtectedRoute>
        <></>
      </ProtectedRoute>
    ),
    children: [
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
];
