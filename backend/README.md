# Yorgut
O.G. Social Media Backend

# TODO
- Authentication
- Communities:
    - "Create Community" page
    - Split communities in two types:
        - Public or Private.
            - Public communities can be joined by anyone. Also, their forums can be viewed by anyone, be it a member or not.
            - Private communities can be joined by anyone, but the owner must accept any new memberships. Also, their forums can only be read by community members.
    - "Community Members" page
        - In the members list of a community, a button "see members" must lead to a page listing all members of given community.
- User:
    - "My Friends" screen
        - Make a screen to list a user's friends. Lacking references, but should be a simple list with a "remove friend" button.
    - "My Communities" screen
        - Same as "My friends"
    - "Scrapbook" page
        - List all user's scraps. Should also contain a text box on top to write a new scrap. The scrapbook can be Public or Private:
            - Public means anyone can read the scraps the user has received.
            - Private means only the receiver can read its scraps.
            - User must be able to toggle between these options anytime.
    - "Search" page
        - Users should be able to search for other users and communities. Probablt not much more than that.
- Profile:
    - "Friends" page
        - In an User Profile, below the friends list, must be a button "see friends" that leads to a full list of the user's friends.
    - "Communities" page
        - In an User Profile, below the communities list, must be a button "see communities" that leads to a full list of the user's communities.
    - "Write Testimonial" page
    - "Testimonials" component
    - "Testimonials" page
- Others:
    - Enable voting in Cool (🧊), Lovable (❤️), Trustworthy (😀)
    - Enable being a fan of someone (probably just like a friendship/membership... fanship?)
    - Add those cool first letters typography
- Add unittests
- Change Sqlite to Postgres for production
