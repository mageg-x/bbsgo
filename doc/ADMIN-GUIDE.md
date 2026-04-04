# RainbowBBS Admin Operation Guide

[中文版](./ADMIN-GUIDE-CN.md)

---

## Table of Contents

- [Foreword](#foreword)
- [Chapter 1: Getting Started with Admin Panel](#chapter-1-getting-started-with-admin-panel)
  - [1.1 Accessing the Admin Panel](#11-accessing-the-admin-panel)
  - [1.2 Logging into the Admin Panel](#12-logging-into-the-admin-panel)
  - [1.3 Admin Panel Interface Overview](#13-admin-panel-interface-overview)
- [Chapter 2: Dashboard](#chapter-2-dashboard)
  - [2.1 Statistics Cards](#21-statistics-cards)
  - [2.2 System Information](#22-system-information)
  - [2.3 Quick Actions](#23-quick-actions)
- [Chapter 3: User Management](#chapter-3-user-management)
  - [3.1 User List](#31-user-list)
  - [3.2 Searching Users](#32-searching-users)
  - [3.3 Editing User Roles](#33-editing-user-roles)
  - [3.4 Muting and Unmuting](#34-muting-and-unmuting)
  - [3.5 Deleting Users](#35-deleting-users)
- [Chapter 4: Post Management](#chapter-4-post-management)
  - [4.1 Post List](#41-post-list)
  - [4.2 Searching Posts](#42-searching-posts)
  - [4.3 Pinning and Unpinning](#43-pinning-and-unpinning)
  - [4.4 Viewing Posts](#44-viewing-posts)
  - [4.5 Deleting Posts](#45-deleting-posts)
- [Chapter 5: Comment Management](#chapter-5-comment-management)
  - [5.1 Comment List](#51-comment-list)
  - [5.2 Searching Comments](#52-searching-comments)
  - [5.3 Viewing Parent Posts](#53-viewing-parent-posts)
  - [5.4 Deleting Comments](#54-deleting-comments)
- [Chapter 6: Board Management](#chapter-6-board-management)
  - [6.1 Board List](#61-board-list)
  - [6.2 Creating Boards](#62-creating-boards)
  - [6.3 Editing Boards](#63-editing-boards)
  - [6.4 Deleting Boards](#64-deleting-boards)
- [Chapter 7: Topic Tag Management](#chapter-7-topic-tag-management)
  - [7.1 Tag List](#71-tag-list)
  - [7.2 Adding Official Tags](#72-adding-official-tags)
  - [7.3 Editing Tags](#73-editing-tags)
  - [7.4 Setting as Official Tags](#74-setting-as-official-tags)
  - [7.5 Disabling and Enabling Tags](#75-disabling-and-enabling-tags)
  - [7.6 Merging Tags](#76-merging-tags)
  - [7.7 Deleting Tags](#77-deleting-tags)
- [Chapter 8: Announcement Management](#chapter-8-announcement-management)
  - [8.1 Announcement List](#81-announcement-list)
  - [8.2 Publishing Announcements](#82-publishing-announcements)
  - [8.3 Editing Announcements](#83-editing-announcements)
  - [8.4 Deleting Announcements](#84-deleting-announcements)
- [Chapter 9: Badge Management](#chapter-9-badge-management)
  - [9.1 Badge List](#91-badge-list)
  - [9.2 Initializing Badges](#92-initializing-badges)
  - [9.3 Viewing Recipients](#93-viewing-recipients)
  - [9.4 Revoking Badges](#94-revoking-badges)
  - [9.5 Deleting Badges](#95-deleting-badges)
- [Chapter 10: Poll Management](#chapter-10-poll-management)
  - [10.1 Poll List](#101-poll-list)
  - [10.2 Searching and Filtering Polls](#102-searching-and-filtering-polls)
  - [10.3 Viewing Poll Details](#103-viewing-poll-details)
  - [10.4 Ending Polls](#104-ending-polls)
  - [10.5 Deleting Polls](#105-deleting-polls)
- [Chapter 11: Report Management](#chapter-11-report-management)
  - [11.1 Report List](#111-report-list)
  - [11.2 Filtering Reports](#112-filtering-reports)
  - [11.3 Processing Reports](#113-processing-reports)
- [Chapter 12: Anti-Spam System Configuration](#chapter-12-anti-spam-system-configuration)
  - [12.1 Rate Limit Configuration](#121-rate-limit-configuration)
  - [12.2 Content Quality Detection Configuration](#122-content-quality-detection-configuration)
  - [12.3 Reputation Score System Configuration](#123-reputation-score-system-configuration)
  - [12.4 Report Processing Configuration](#124-report-processing-configuration)
  - [12.5 User Reputation Management](#125-user-reputation-management)
- [Chapter 13: System Settings](#chapter-13-system-settings)
  - [13.1 Feature Switches](#131-feature-switches)
  - [13.2 Point Rules](#132-point-rules)
  - [13.3 Saving and Resetting](#133-saving-and-resetting)
- [Chapter 14: Main Site User Guide](#chapter-14-main-site-user-guide)
  - [14.1 Main Site Overview](#141-main-site-overview)
  - [14.2 User Registration and Login](#142-user-registration-and-login)
  - [14.3 Browsing and Searching](#143-browsing-and-searching)
  - [14.4 Posting and Commenting](#144-posting-and-commenting)
  - [14.5 Social Interactions](#145-social-interactions)
  - [14.6 User Profile Center](#146-user-profile-center)
- [Chapter 15: FAQ and Troubleshooting](#chapter-15-faq-and-troubleshooting)
  - [15.1 Admin Panel Issues](#151-admin-panel-issues)
  - [15.2 Main Site Issues](#152-main-site-issues)
  - [15.3 Database Issues](#153-database-issues)
- [Chapter 16: Best Practices](#chapter-16-best-practices)
  - [16.1 Daily Operation Recommendations](#161-daily-operation-recommendations)
  - [16.2 Community Management Strategies](#162-community-management-strategies)
  - [16.3 Security Recommendations](#163-security-recommendations)

---

## Foreword

Welcome to the RainbowBBS Admin Operation Guide! This document provides detailed information about all features of the RainbowBBS community forum system's admin panel, as well as how to use the main site.

### Intended Audience

- Community administrators
- Website operators
- Technical maintenance personnel

### Prerequisites

Before reading this document, please ensure:
1. You have successfully deployed the RainbowBBS system
2. You have an administrator account (default: admin/12345678)
3. You have basic computer operation knowledge

### Terminology

| Term | Description |
|------|-------------|
| Board | Discussion area categories, such as Technology Exchange, Q&A, etc. |
| Topic Tag | Post classification tags, customizable by users |
| Official Tag | Administrator-recommended tags, displayed with priority |
| Badge | User achievement identifiers, automatically awarded when conditions are met |
| Reputation Score | User credit score, affects content visibility |
| Pinned | Posts fixed at the top of the list |
| Featured | Posts marked as high-quality content |

---

## Chapter 1: Getting Started with Admin Panel

### 1.1 Accessing the Admin Panel

The admin panel is typically deployed on a separate port or path.

#### Access URL

Depending on your deployment, the access URL may be:
- `http://your-domain.com/admin`
- `http://your-domain.com:8081` (if using a separate port)

#### Local Development Environment

For local development:
- Main site typically runs at: `http://localhost:5173`
- Admin panel typically runs at: `http://localhost:5174` (or another port)

### 1.2 Logging into the Admin Panel

#### Login Steps

1. Open the admin panel access URL
2. Enter on the login page:
   - **Username**: admin (or your administrator username)
   - **Password**: 12345678 (default password, please change immediately after first login)
3. Click the "Login" button

#### Default Account

| Account Type | Username | Password | Description |
|--------------|----------|----------|-------------|
| Administrator | admin | 12345678 | Full access |

#### Changing Password

Please change the administrator password immediately after first login:

1. Login to the admin panel
2. Find the "Change Password" menu (usually in the user avatar dropdown)
3. Enter current password and new password
4. Click "Confirm" to save

#### Security Tips

⚠️ **Important Security Tips**:
- Do NOT use the default password in production
- Change the administrator password regularly
- Do NOT share the administrator account password with others
- Use strong passwords (include uppercase, lowercase, numbers, and special characters)

### 1.3 Admin Panel Interface Overview

#### Left Navigation Bar

The admin panel uses a left-side navigation layout. Main menu items include:

| Menu Item | Function Description |
|-----------|----------------------|
| Dashboard | View statistics and system information |
| User Management | Manage user accounts, roles, mutes, etc. |
| Post Management | Manage all posts, pinning, deletion, etc. |
| Comment Management | Manage comment content |
| Board Management | Manage discussion boards |
| Tag Management | Manage topic tags |
| Announcement Management | Publish and manage announcements |
| Badge Management | Manage user badges |
| Poll Management | Manage poll activities |
| Report Management | Handle user reports |
| Anti-Spam Config | Configure anti-spam system parameters |
| System Settings | Configure system feature switches and point rules |

#### Top Bar

- **Left**: System logo and name
- **Right**: Current user information, logout button

#### Content Area

- Displays the function page of the currently selected menu
- Contains tables, forms, buttons, and other operation elements

---

## Chapter 2: Dashboard

The dashboard is the home page of the admin panel, displaying key system statistics and quick action shortcuts.

### 2.1 Statistics Cards

The top of the dashboard displays 4 key statistics cards:

#### User Statistics

- **Display**: Total users
- **Trend Arrow**: Shows percentage increase/decrease compared to previous period
- **Color**: Purple gradient
- **Icon**: User icon

#### Post Statistics

- **Display**: Total posts
- **Trend Arrow**: Shows percentage increase/decrease compared to previous period
- **Color**: Green gradient
- **Icon**: File icon

#### Comment Statistics

- **Display**: Total comments
- **Trend Arrow**: Shows percentage increase/decrease compared to previous period
- **Color**: Cyan gradient
- **Icon**: Message icon

#### Pending Reports Statistics

- **Display**: Number of pending reports
- **Trend Arrow**: Shows percentage increase/decrease compared to previous period
- **Color**: Red gradient
- **Icon**: Warning triangle icon

### 2.2 System Information

Below the statistics cards is the system information card, containing:

#### Version Information

- **System Version**: v1.0.0
- **Icon**: Code icon
- **Color**: Blue

#### Go Version

- **Go Version**: go1.21+
- **Icon**: Box icon
- **Color**: Green

#### Database

- **Database Type**: SQLite
- **Icon**: Database icon
- **Color**: Purple

#### Cache

- **Cache System**: Ristretto
- **Icon**: Lightning icon
- **Color**: Orange

#### Uptime

- **System Uptime**: Shows days and hours the system has been running
- **Icon**: CPU icon
- **Color**: Cyan

#### System Status

- **Status**: Normal (Online)
- **Icon**: Activity icon
- **Color**: Pink

### 2.3 Quick Actions

Below the system information is the quick actions area, providing quick access to 4 commonly used functions:

#### User Management Quick Access

- **Icon**: User plus icon
- **Color**: Purple
- **Function**: Click to jump directly to User Management page

#### Post Moderation Quick Access

- **Icon**: File plus icon
- **Color**: Green
- **Function**: Click to jump directly to Post Management page

#### Handle Reports Quick Access

- **Icon**: Shield warning icon
- **Color**: Red
- **Function**: Click to jump directly to Report Management page

#### Publish Announcement Quick Access

- **Icon**: Megaphone icon
- **Color**: Yellow
- **Function**: Click to jump directly to Announcement Management page

---

## Chapter 3: User Management

The User Management module manages all user accounts in the community.

### 3.1 User List

After entering the "User Management" menu, the user list table is displayed:

#### Table Columns

| Column | Description |
|--------|-------------|
| ID | User unique identifier |
| User Info | Displays user avatar, username, email |
| Role | User role (Regular User/Moderator/Administrator) |
| Points | User points, sortable |
| Status | Normal/Muted |
| Registration Time | User registration date |
| Actions | Role edit, Mute/Unmute, Delete |

#### User Role Description

| Role | ID | Color | Permission Description |
|------|----|-------|-----------------------|
| Regular User | 0 | Gray | Basic posting and commenting permissions |
| Moderator | 1 | Orange | Can manage specified board content |
| Administrator | 2 | Red | Full admin permissions |

### 3.2 Searching Users

#### Search Steps

1. Enter keywords in the search box at the top right of the page
2. Keywords can be username or email
3. Click the "Search" button or press Enter to execute the search
4. Search results update the table in real-time

#### Clearing Search

- Click the "×" button on the right side of the search box to clear search keywords
- Or directly delete the content in the search box and press Enter

### 3.3 Editing User Roles

#### Modifying Role Steps

1. Find the user to modify in the user list
2. Click the "Role" button in that row
3. Select a new role in the popup dialog:
   - Regular User
   - Moderator
   - Administrator
4. Click the "Confirm" button to save

#### Notes

⚠️ **Important Notes**:
- Only administrators can modify user roles
- Role changes take effect immediately
- Be cautious when granting administrator privileges

### 3.4 Muting and Unmuting

#### Muting a User

1. Find the user to mute in the user list
2. Click the "Mute" button
3. Click "OK" in the confirmation dialog
4. User status changes to "Muted"

#### Unmuting a User

1. Find the muted user in the user list
2. Click the "Unmute" button
3. Click "OK" in the confirmation dialog
4. User status returns to "Normal"

#### Mute Impact

Muted users cannot:
- Post new threads
- Post comments
- Send private messages

But can still:
- Browse posts
- Like and bookmark
- View their profile

### 3.5 Deleting Users

#### Deletion Steps

⚠️ **Warning**: Deleting a user is irreversible, please proceed with caution!

1. Find the user to delete in the user list
2. Click the "Delete" button
3. Carefully read the prompt in the confirmation dialog:
   > Are you sure you want to permanently delete user "xxx"? This action cannot be undone, and related posts, bookmarks, likes, and other data will also be deleted!
4. Click the "Delete" button after confirming

#### Deletion Consequences

Deleting a user will also delete:
- User account
- All posts by the user
- All comments by the user
- User's bookmark records
- User's like records
- User's follow relationships
- User's earned badges

---

## Chapter 4: Post Management

The Post Management module manages all post content in the community.

### 4.1 Post List

After entering the "Post Management" menu, the post list table is displayed:

#### Table Columns

| Column | Description |
|--------|-------------|
| ID | Post unique identifier |
| Title | Post title, clickable to open directly in new window |
| Author | Avatar and username of the posting user |
| Board | Board the post belongs to |
| Views/Likes/Comments | Post views, likes, comments count |
| Status | Pinned, Featured, Locked indicators |
| Published Time | Post publication date |
| Actions | Pin/Unpin, View, Delete |

#### Status Tag Description

| Tag | Color | Description |
|-----|-------|-------------|
| Pinned | Red | Post fixed at top of list |
| Featured | Yellow | Marked as high-quality content |
| Locked | Gray | Comments disabled |

### 4.2 Searching Posts

#### Search Steps

1. Enter keywords in the search box at the top right of the page
2. Keywords match post titles
3. Click the "Search" button or press Enter to execute the search

#### Clearing Search

- Click the "×" button on the right side of the search box to clear search

### 4.3 Pinning and Unpinning

#### Pinning a Post

1. Find the post to pin in the post list
2. Click the "Pin" button
3. Click "OK" in the confirmation dialog
4. Post status shows "Pinned" tag

#### Unpinning a Post

1. Find the pinned post in the post list
2. Click the "Unpin" button
3. Click "OK" in the confirmation dialog
4. Post "Pinned" tag disappears

#### Pinning Rules

- Pinned posts appear at the top of the board list
- Multiple posts can be pinned simultaneously
- Pinned posts are sorted by pin time in reverse order

### 4.4 Viewing Posts

#### Viewing Steps

1. Find the post to view in the post list
2. Click the "View" button, or directly click the post title
3. Post opens in a new tab

### 4.5 Deleting Posts

#### Deletion Steps

⚠️ **Warning**: Deleting a post is irreversible!

1. Find the post to delete in the post list
2. Click the "Delete" button
3. Click "Delete" in the confirmation dialog
4. Post is removed from the list

#### Deletion Consequences

Deleting a post will also delete:
- Post content
- All comments on the post
- Post's like records
- Post's bookmark records
- Post's poll (if any)

---

## Chapter 5: Comment Management

The Comment Management module manages all comment content in the community.

### 5.1 Comment List

After entering the "Comment Management" menu, the comment list table is displayed:

#### Table Columns

| Column | Description |
|--------|-------------|
| ID | Comment unique identifier |
| Content | Comment content, long content is truncated |
| Author | Avatar and username of the commenting user |
| Parent Post | Click "View" button to jump to original post |
| Likes | Number of likes the comment received |
| Published Time | Comment publication date |
| Actions | Delete |

### 5.2 Searching Comments

#### Search Steps

1. Enter keywords in the search box at the top right of the page
2. Keywords match comment content
3. Click the "Search" button or press Enter to execute the search

#### Viewing Full Content

- Hover over comment content to see a tooltip with the full content

### 5.3 Viewing Parent Posts

#### Viewing Steps

1. Find the comment to view in the comment list
2. Click the "View" button
3. Original post opens in a new tab and automatically scrolls to the comment location

### 5.4 Deleting Comments

#### Deletion Steps

⚠️ **Warning**: Deleting a comment is irreversible!

1. Find the comment to delete in the comment list
2. Click the "Delete" button
3. Click "Delete" in the confirmation dialog
4. Comment is removed from the list

#### Deletion Consequences

Deleting a comment will also delete:
- Comment content
- Comment's like records
- Replies to this comment (if any)

---

## Chapter 6: Board Management

The Board Management module manages the community's discussion boards.

### 6.1 Board List

After entering the "Board Management" menu, the board list table is displayed:

#### Table Columns

| Column | Description |
|--------|-------------|
| ID | Board unique identifier |
| Name | Board name with icon |
| Description | Board description |
| Sort | Sort value, sortable |
| Allow Posting | Whether users can post in this board |
| Actions | Edit, Delete |

#### Default Boards

System comes with 8 preset boards:

| Board Name | Description |
|------------|-------------|
| All | Shows posts from all boards (virtual board) |
| Technology | Technology discussion area |
| Q&A | Question and answer area |
| News | News and information area |
| Resources | Resource sharing area |
| Jobs | Recruitment and job search area |
| Off-Topic | Casual chat area |
| Site Feedback | Site management discussion area |

### 6.2 Creating Boards

#### Creation Steps

1. Click the "Create Board" button at the top right of the page
2. Fill in the information in the popup dialog:

| Field | Required | Description |
|-------|----------|-------------|
| Board Name | Yes | Board display name |
| Board Description | No | Detailed board description |
| Sort | No | Sort value, smaller appears first |
| Allow Posting | No | Check to allow users to post in this board |

3. Click the "Save" button

#### Sorting Rules

- Smaller sort values make the board appear earlier in the navigation
- Same sort values are ordered by creation time
- Recommended to use intervals like 0, 10, 20, 30... to make inserting new boards easier later

### 6.3 Editing Boards

#### Editing Steps

1. Find the board to edit in the board list
2. Click the "Edit" button
3. Modify information in the popup dialog
4. Click the "Save" button

### 6.4 Deleting Boards

#### Deletion Steps

⚠️ **Warning**: Ensure there are no posts in this board before deleting!

1. Find the board to delete in the board list
2. Click the "Delete" button
3. Click "Delete" in the confirmation dialog
4. Board is removed from the list

#### Notes

⚠️ **Important Notes**:
- Deleting a board does not automatically delete posts in the board
- After deleting a board, posts in that board will not display properly
- Recommend moving posts to another board first before deleting

---

## Chapter 7: Topic Tag Management

The Topic Tag Management module manages tags used when users post.

### 7.1 Tag List

After entering the "Topic Tag Management" menu, the tag list table is displayed:

#### Table Columns

| Column | Description |
|--------|-------------|
| ID | Tag unique identifier |
| Name | Tag name with icon (if any) |
| Usage Count | Number of times this tag has been used, sortable |
| Status | Normal/Disabled |
| Created Time | Tag creation date |
| Actions | Edit, Set as Official, Disable/Enable, Merge, Delete |

#### Tag Types

| Type | Indicator | Description |
|------|-----------|-------------|
| Official Tag | Blue tag | Administrator recommended, priority display |
| Regular Tag | No indicator | User-created tag |
| Disabled | Strikethrough | Unusable tag |

#### Prompt Information

There's a prompt at the top of the page:
> Tags are automatically created when users post, administrators only need to handle problematic tags

This means:
- New tags entered by users when posting are automatically created
- Administrators don't need to manually create every tag
- Administrators only need to handle problematic tags (disable, delete, etc.)

### 7.2 Adding Official Tags

Although users can automatically create tags, administrators can add official recommended tags.

#### Adding Steps

1. Click the "Add Official Tag" button at the top right of the page
2. Fill in the information in the popup dialog:

| Field | Required | Description |
|-------|----------|-------------|
| Name | Yes | Tag name |
| Icon (Emoji) | No | Tag's Emoji icon, e.g., 🌐 |
| Set as Official Recommended Tag | No | Check to mark as official tag |

3. Click the "Create" button

#### Icon Preview

- The entered Emoji is previewed in real-time to the right of the icon input box
- If no icon is entered, the default 👁 icon is displayed

### 7.3 Editing Tags

#### Editing Steps

1. Find the tag to edit in the tag list
2. Click the "Edit" button
3. Modify information in the popup dialog:
   - Tag name
   - Tag icon
4. Click the "Update" button

### 7.4 Setting as Official Tags

#### Setting as Official

1. Find the tag to set in the tag list
2. Click the "Set as Official" button
3. Tag shows "Official" indicator

#### Unsetting as Official

1. Find the tag already set as official in the tag list
2. Click the "Unset as Official" button
3. Tag's "Official" indicator disappears

#### Official Tag Advantages

- Official tags are displayed with priority in the tag list
- Official tags are recommended first when users post
- Official tags have an "Official" indicator for more authority

### 7.5 Disabling and Enabling Tags

#### Disabling Tags

1. Find the tag to disable in the tag list
2. Click the "Disable" button
3. Click "OK" in the confirmation dialog
4. Tag status changes to "Disabled", name shows strikethrough

#### Enabling Tags

1. Find the disabled tag in the tag list
2. Click the "Enable" button
3. Click "OK" in the confirmation dialog
4. Tag status returns to "Normal"

#### Disable Impact

- Disabled tags cannot be used when posting
- Posts already using this tag still show the tag
- But new posts can no longer select this tag

### 7.6 Merging Tags

When similar or duplicate tags exist, they can be merged.

#### Merging Steps

1. Find the source tag to merge (tag to be deleted) in the tag list
2. Click the "Merge" button
3. In the popup dialog:
   - Confirm the source tag name
   - Select the target tag (tag to keep) from the dropdown
4. Click the "Confirm Merge" button

#### Merge Effect

- All usage records from source tag are transferred to target tag
- Source tag is deleted
- Target tag's usage count increases

#### Notes

⚠️ **Important Notes**:
- Merge operations are irreversible
- Be careful when selecting the target tag
- Recommend backing up data first

### 7.7 Deleting Tags

#### Deletion Steps

⚠️ **Warning**: Deleting a tag is irreversible!

1. Find the tag to delete in the tag list
2. Click the "Delete" button
3. Click "Delete" in the confirmation dialog
4. Tag is removed from the list

#### Deletion Consequences

- Tag is permanently deleted
- Posts that used this tag no longer show the tag
- Tag's usage count statistics are lost

---

## Chapter 8: Announcement Management

The Announcement Management module publishes and manages community announcements.

### 8.1 Announcement List

After entering the "Announcement Management" menu, the announcement list table is displayed:

#### Table Columns

| Column | Description |
|--------|-------------|
| ID | Announcement unique identifier |
| Title | Announcement title, pinned announcements have an icon indicator |
| Pinned | Whether pinned |
| Expiration Time | Announcement expiration time, "Permanent" if never expires |
| Published Time | Announcement publication date |
| Actions | Edit, Delete |

### 8.2 Publishing Announcements

#### Publishing Steps

1. Click the "Publish Announcement" button at the top right of the page
2. Fill in the information in the popup dialog:

| Field | Required | Description |
|-------|----------|-------------|
| Title | Yes | Announcement title |
| Content | Yes | Detailed announcement content |
| Pin Announcement | No | Check to pin the announcement |
| Expiration Time (Optional) | No | Set announcement expiration time, leave blank for permanent |

3. Click the "Save" button

#### Pinned Announcements

- Pinned announcements appear at the top of all announcements
- Multiple announcements can be pinned simultaneously
- Pinned announcements have a pushpin icon indicator

#### Expiration Time

- No expiration time set: Announcement displays permanently
- Expiration time set: Announcement automatically hides after reaching the specified time
- Expired announcements are still visible in the admin panel but not on the front-end

### 8.3 Editing Announcements

#### Editing Steps

1. Find the announcement to edit in the announcement list
2. Click the "Edit" button
3. Modify information in the popup dialog
4. Click the "Save" button

### 8.4 Deleting Announcements

#### Deletion Steps

1. Find the announcement to delete in the announcement list
2. Click the "Delete" button
3. Click "Delete" in the confirmation dialog
4. Announcement is removed from the list

---

## Chapter 9: Badge Management

The Badge Management module manages user achievement badges.

### 9.1 Badge List

After entering the "Badge Management" menu, the badge list table is displayed:

#### Table Columns

| Column | Description |
|--------|-------------|
| ID | Badge unique identifier |
| Badge | Badge icon, name, type |
| Description | Detailed badge description |
| Requirements | Requirements to earn this badge |
| Recipient Count | Number of users who have earned this badge |
| Actions | View Users, Delete |

#### Badge Types

| Type | Description |
|------|-------------|
| Beginner | Easy-to-earn basic badges |
| Intermediate | Badges requiring some effort |
| Advanced | Very rare high-level badges |

#### Default Badge List

System comes with 9 preset badges:

| Badge Name | Type | Requirements |
|------------|------|--------------|
| Newcomer | Beginner | Earned upon registration |
| First Post | Beginner | Publish 1st post |
| First Comment | Beginner | Publish 1st comment |
| Prolific Author | Intermediate | Publish 50 posts |
| Chatty Cathy | Intermediate | Publish 1000 comments |
| Popular | Intermediate | Get 1000 cumulative likes |
| Top Commenter | Intermediate | Get 5 best comments |
| Influencer | Advanced | Followed by 500 users |
| Community Legend | Advanced | Registered ≥ 2 years + Posts ≥ 200 + Likes ≥ 500 + Best comments ≥ 10 |

### 9.2 Initializing Badges

If badge data is lost or needs reset, system badges can be initialized.

#### Initialization Steps

⚠️ **Warning**: Initialization will reset all badge data!

1. Click the "Initialize Badges" button at the top right of the page
2. Read the prompt in the confirmation dialog:
   > Are you sure you want to initialize system badges? This will create 9 default badges.
3. Click the "OK" button
4. System badge initialization is complete

### 9.3 Viewing Recipients

#### Viewing Steps

1. Find the badge to view in the badge list
2. Click the "View Users" button
3. List of users who earned this badge is displayed in the popup dialog

#### User List Table

| Column | Description |
|--------|-------------|
| ID | User unique identifier |
| User Info | User avatar, username, email |
| Earned Time | Time user earned this badge |
| Status | Normal/Revoked |
| Actions | Revoke (only for normal status) |

#### Pagination

- If more than 10 recipients, pagination controls are displayed
- Click page numbers to switch pages

### 9.4 Revoking Badges

If a user violates rules or earned a badge in error, it can be revoked.

#### Revocation Steps

1. Find the user to revoke in the badge's recipient list
2. Click the "Revoke" button
3. Fill in the revocation reason in the popup dialog (required)
4. Click the "Confirm Revocation" button

#### Revocation Reason

Revocation reason is required, please detail:
- Violation behavior
- System error
- Other reasons

#### Revocation Impact

- User loses the badge
- Badge recipient count decreases by 1
- User receives notification that badge was revoked

### 9.5 Deleting Badges

#### Deletion Steps

⚠️ **Warning**: Deleting a badge is irreversible!

1. Find the badge to delete in the badge list
2. Click the "Delete" button
3. Click "Delete" in the confirmation dialog
4. Badge is removed from the list

#### Deletion Consequences

- Badge is permanently deleted
- All users lose the badge
- Badge earning records are lost

---

## Chapter 10: Poll Management

The Poll Management module manages poll activities in posts.

### 10.1 Poll List

After entering the "Poll Management" menu, the poll list table is displayed:

#### Table Columns

| Column | Description |
|--------|-------------|
| ID | Poll unique identifier |
| Poll Title | Poll title, clickable to jump to original post |
| Type | Single choice/Multiple choice |
| Option Count | Number of poll options |
| Voter Count | Number of participating users |
| Status | Active/Ended |
| Deadline | Poll deadline, "Permanent" if never expires |
| Created Time | Poll creation time |
| Actions | View, End, Delete |

#### Poll Types

| Type | Indicator | Description |
|------|-----------|-------------|
| Single Choice | Blue tag | User can only select one option |
| Multiple Choice | Green tag | User can select multiple options |

#### Poll Status

| Status | Description |
|--------|-------------|
| Active | Users can vote |
| Ended | Users cannot vote, can only view results |

### 10.2 Searching and Filtering Polls

#### Searching Polls

1. Enter keywords in the search box at the top right of the page
2. Keywords match poll titles
3. Click the "Search" button or press Enter to execute the search

#### Status Filter

1. Click the "Status Filter" dropdown
2. Select filter condition:
   - All statuses
   - Active
   - Ended
3. List automatically refreshes

### 10.3 Viewing Poll Details

#### Viewing Steps

1. Find the poll to view in the poll list
2. Click the "View" button
3. Poll details are displayed in the popup dialog

#### Details Content

- **Poll Title**: Poll title
- **Poll Type**: Single or multiple choice (multiple shows max selectable)
- **Participant Count**: Total voters
- **Option List**:
  - Text of each option
  - Vote count and percentage for each option
  - Progress bar showing vote ratio

### 10.4 Ending Polls

If a poll needs to end early, it can be ended manually.

#### Ending Steps

⚠️ **Note**: Only active polls can be ended!

1. Find the poll to end in the poll list
2. Click the "End" button
3. Click "OK" in the confirmation dialog
4. Poll status changes to "Ended"

#### Ending Impact

- Users can no longer vote
- Users can still view poll results
- Poll cannot be re-opened

### 10.5 Deleting Polls

#### Deletion Steps

⚠️ **Warning**: Deleting a poll is irreversible!

1. Find the poll to delete in the poll list
2. Click the "Delete" button
3. Click "Delete" in the confirmation dialog
4. Poll is removed from the list

#### Deletion Consequences

- Poll is permanently deleted
- All vote records are lost
- Poll also disappears from the post

---

## Chapter 11: Report Management

The Report Management module handles reports submitted by users.

### 11.1 Report List

After entering the "Report Management" menu, the report list table is displayed:

#### Table Columns

| Column | Description |
|--------|-------------|
| ID | Report unique identifier |
| Reporter | User who submitted the report |
| Reported User | Author of reported content (if any) |
| Type | Reported content type (Post/Comment/Private Message) |
| Reported Content | Summary of reported content, clickable to view |
| Report Reason | Type of report reason |
| Notes | Additional notes from reporter |
| Status | Pending/Approved/Rejected |
| Report Time | Report submission time |
| Actions | Approve/Reject (only for pending) |

#### Pending Indicator

Top left of page shows pending report count:
> 📢 X pending

Red indicator reminds administrators to process promptly.

#### Report Types

| Type | Indicator | Description |
|------|-----------|-------------|
| Post | Blue tag | Report is about a post |
| Comment | Blue tag | Report is about a comment |
| Private Message | Orange tag | Report is about a private message |

#### Report Reasons

| Reason | Description |
|--------|-------------|
| Spam/Ad | Content is advertising or spam |
| Inappropriate Content | Content violates community rules |
| Harassment | Content contains personal attacks |
| Misinformation | Content is rumors or false information |
| Other | Other reasons |

#### Report Status

| Status | Color | Description |
|--------|-------|-------------|
| Pending | Orange | Needs administrator action |
| Approved | Green | Report valid, content deleted |
| Rejected | Gray | Report invalid |

### 11.2 Filtering Reports

#### Status Filter

1. Click the "Select Status" dropdown
2. Select filter condition:
   - All statuses
   - Pending
   - Approved
   - Rejected
3. List automatically refreshes

#### Refreshing List

Click the "Refresh" button to reload the report list.

### 11.3 Processing Reports

#### Viewing Reported Content

- Click "View Post" or "View Comment" link
- Reported content opens in new tab
- Carefully verify if report is accurate

#### Approving Reports

If report is valid:

1. Click the "Approve" button
2. Read the prompt in the confirmation dialog:
   > After approval, related content will be deleted. Continue?
3. Click the "Approve" button
4. Report status changes to "Approved"
5. Reported content is deleted

#### Rejecting Reports

If report is not valid:

1. Click the "Reject" button
2. Click "OK" in the confirmation dialog
3. Report status changes to "Rejected"
4. Reported content remains unchanged

#### Processing Consequences

**Approving Report**:
- Reported content is deleted
- Reported user's reputation score decreases
- Reporter gets reputation score bonus

**Rejecting Report**:
- Reported content remains
- Reporter's reputation score may decrease (malicious report)

---

## Chapter 12: Anti-Spam System Configuration

The Anti-Spam System Configuration module configures the community's anti-spam strategies and parameters.

### 12.1 Rate Limit Configuration

Click the "Rate Limits" tab to configure operation rate limits.

#### Configuration Items

| Config Item | Default | Range | Description |
|-------------|---------|-------|-------------|
| Min post interval (seconds) | 60 | 10-600 | Minimum time between user posts |
| Min comment interval (seconds) | 30 | 5-300 | Minimum time between user comments |
| Max posts per day | 10 | 1-100 | Maximum posts per user per day |
| Max comments per day | 50 | 1-200 | Maximum comments per user per day |

#### New User Limits

Additional limits for new users:

| Config Item | Default | Range | Description |
|-------------|---------|-------|-------------|
| New user duration (hours) | 24 | 1-168 | How long after registration considered new user |
| New user max posts per day | 3 | 1-20 | Maximum posts per new user per day |
| New user max comments per day | 10 | 1-50 | Maximum comments per new user per day |

#### Saving Configuration

After modifying configuration, click the "Save Config" button to save.

#### Configuration Recommendations

**Small Communities**:
- Post interval: 30-60 seconds
- Daily posts: 10-20 posts
- New user limits: Relatively lenient

**Medium Communities**:
- Post interval: 60-120 seconds
- Daily posts: 5-10 posts
- New user limits: Moderate

**Large Communities**:
- Post interval: 120-300 seconds
- Daily posts: 3-5 posts
- New user limits: Strict

### 12.2 Content Quality Detection Configuration

Click the "Content Quality" tab to configure content quality detection parameters.

#### Configuration Items

| Config Item | Default | Range | Description |
|-------------|---------|-------|-------------|
| Min content length | 10 | 5-100 | Minimum Chinese characters for post content |
| Repeated character threshold | 5 | 3-20 | Consecutive repeated characters beyond this considered spam |
| Content similarity threshold | 0.8 | 0.5-1.0 | Beyond this considered duplicate content |

#### Saving Configuration

After modifying configuration, click the "Save Config" button to save.

### 12.3 Reputation Score System Configuration

Click the "Reputation System" tab to configure reputation score related parameters.

#### Reputation Score Thresholds

| Config Item | Default | Range | Description |
|-------------|---------|-------|-------------|
| Low reputation threshold | 60 | 20-80 | Below this score content heat decreases |
| Mute reputation threshold | 20 | 0-50 | Below this score automatically muted |
| Enable low reputation mute | On | - | When enabled, low reputation users are automatically muted |

#### Heat Weights

| Config Item | Default | Range | Description |
|-------------|---------|-------|-------------|
| Low quality content heat coefficient | 0.3 | 0-1 | Heat multiplier for low quality content |
| Low reputation content heat coefficient | 0.5 | 0-1 | Heat multiplier for low reputation user content |

#### Saving Configuration

After modifying configuration, click the "Save Config" button to save.

#### Reputation Score Levels

| Reputation | Level | Color | Description |
|------------|-------|-------|-------------|
| 80-100 | Normal | Green | Content displays normally |
| 60-79 | Needs Verification | Blue | Content heat slightly reduced |
| 40-59 | Restricted | Orange | Content heat reduced |
| 20-39 | Severe | Red | Content heat significantly reduced |
| 0-19 | Muted | Gray | Automatically muted |

### 12.4 Report Processing Configuration

Click the "Report Handling" tab to configure automatic report handling rules.

#### Configuration Items

| Config Item | Default | Range | Description |
|-------------|---------|-------|-------------|
| Auto-hide report threshold | 3 | 1-10 | Auto-hide content when reaching this report count |
| Auto-mute violation threshold | 5 | 1-20 | Auto-mute when hidden content reaches this in 7 days |
| Auto-mute days | 3 | 1-30 | Days for auto-mute |
| Daily report limit | 10 | 1-50 | Maximum reports per user per day |

#### Saving Configuration

After modifying configuration, click the "Save Config" button to save.

### 12.5 User Reputation Management

Click the "User Reputation" tab to manage user reputation scores.

#### Searching Users

1. Enter username or nickname in the search box
2. List filters in real-time to show matching users

#### User List

| Column | Description |
|--------|-------------|
| ID | User unique identifier |
| Username | User login name |
| Nickname | User display name |
| Reputation Score | Reputation score progress bar, color indicates level |
| Reputation Level | Normal/Needs Verification/Restricted/Severe/Muted |
| Registration Time | User registration time |
| Actions | Adjust Reputation, Mute/Unmute |

#### Adjusting Reputation Score

1. Find the user to adjust in the user list
2. Click the "Adjust Reputation" button
3. In the popup dialog:
   - View current reputation score
   - Enter adjustment points (positive to add, negative to subtract)
   - Fill in adjustment reason
4. Click the "Confirm" button

#### Manual Mute/Unmute

**Muting**:
1. Find an unmuted user
2. Click the "Mute" button
3. Click "OK" in the confirmation dialog

**Unmuting**:
1. Find a muted user
2. Click the "Unmute" button
3. User is automatically unmuted

#### Pagination

- 20 users per page
- Click page numbers to switch pages

---

## Chapter 13: System Settings

The System Settings module configures the community's basic features and point rules.

### 13.1 Feature Switches

#### Switch List

| Switch | Default | Description |
|--------|---------|-------------|
| Allow User Registration | On | Control whether new user registration is open |
| Allow Posting | On | Control whether users can publish posts |
| Allow Commenting | On | Control whether users can post comments |
| Allow Poll Creation | On | Control whether users can publish polls |

#### Switch Descriptions

**Allow User Registration**:
- On: New users can register accounts
- Off: Stop new user registration, only existing users can login

**Allow Posting**:
- On: Users can publish new posts
- Off: All users cannot post (except administrators)

**Allow Commenting**:
- On: Users can post comments
- Off: All users cannot comment (except administrators)

**Allow Poll Creation**:
- On: Users can add polls to posts
- Off: Users cannot add polls

#### Usage Scenarios

**Community Maintenance**:
- Temporarily close registration, only allow existing users access
- Temporarily close posting to clean up inappropriate content

**Special Events**:
- Only allow comments, not new posts
- Restrict certain features

### 13.2 Point Rules

#### Point Configuration Items

| Config Item | Default | Range | Description |
|-------------|---------|-------|-------------|
| Post points | 5 | 0-999 | Points earned for publishing a post |
| Comment points | 1 | 0-999 | Points earned for posting a comment |
| Check-in points | 2 | 0-999 | Points earned for daily check-in |

#### Point Rule Descriptions

**Post points**:
- User earns these points each time they publish a new post
- Points are not deducted if post is deleted
- Recommended: 5-10 points

**Comment points**:
- User earns these points each time they post a comment
- Points are not deducted if comment is deleted
- Recommended: 1-3 points

**Check-in points**:
- User earns these points for daily check-in
- Consecutive check-ins earn extra bonus (configurable in code)
- Recommended: 2-5 points

#### Point Uses

Points can be used for:
- Unlocking special features
- Redeeming virtual items
- Improving user level
- Displaying user achievements

### 13.3 Saving and Resetting

#### Saving Settings

1. After modifying configuration
2. Click the "Save Settings" button
3. "Settings saved" prompt appears
4. Configuration takes effect immediately

#### Resetting Settings

1. Click the "Reset" button
2. Configuration restores to last saved values
3. Does not restore to system defaults

#### Notes

⚠️ **Important Notes**:
- Modifying feature switches immediately affects all users
- Recommend modifying during low traffic periods
- Notify users before modifying
- Observe community reaction after modification

---

## Chapter 14: Main Site User Guide

This chapter introduces how to use the main site (user side), helping administrators understand how users use the community.

### 14.1 Main Site Overview

#### Page Layout

- **Top Navigation Bar**:
  - Logo and site name
  - Search button
  - Notification button
  - Message button
  - User menu (after login)

- **Left Sidebar**:
  - Hot topic tags

- **Main Content Area**:
  - Announcements (if any)
  - Post list

- **Right Sidebar**:
  - Check-in button
  - Hot posts
  - Active users

#### Board Navigation

Board tabs are displayed at the top:
- All
- Technology
- Q&A
- News
- Resources
- Jobs
- Off-Topic
- Site Feedback

Click a tab to switch to that board.

### 14.2 User Registration and Login

#### Registering an Account

1. Click the "Register" button at top right
2. Fill in the registration form:
   - Username: 3-20 characters
   - Email: For password recovery
   - Password: 6-20 characters
3. Click the "Register" button
4. Auto-login after successful registration
5. Earn "Newcomer" badge

#### Logging In

1. Click the "Login" button at top right
2. Enter username and password
3. Click the "Login" button

#### Logging Out

1. Click the user avatar at top right
2. Select "Logout"

### 14.3 Browsing and Searching

#### Browsing Boards

1. Click a board tab at the top
2. Browse that board's post list

#### Viewing Posts

1. Click any post card
2. Enter post detail page
3. View:
   - Post content (Markdown rendered)
   - Images (click to enlarge)
   - Poll (if any)
   - Comment list

#### Searching Posts

1. Click the "Search" button at the top
2. Enter keywords
3. Press Enter or click search
4. View search results

### 14.4 Posting and Commenting

#### Publishing a New Post

1. Click the "Post" button at the top
2. Fill in post information:
   - Title (required)
   - Select board (required)
   - Content (use Markdown editor)
   - Tags (up to 3, optional)
   - Poll (optional)

3. Click the "Publish" button

#### Markdown Editor

**Features**:
- Supports GFM syntax
- Supports code highlighting
- Supports math formulas
- Supports image upload
- Supports video upload

**Inserting Images**:
- Drag and drop images into editor
- Copy and paste images
- On small screens click "Insert Image" button

**Creating Polls**:
1. Click "Add Poll"
2. Fill in poll title
3. Add options (2-10)
4. Select single/multiple choice
5. Set deadline (optional)
6. Click "Confirm"

#### Posting Comments

1. At the bottom of the post detail page
2. Enter comment content
3. Click "Post Comment"

#### Replying to Comments

1. Click "Reply" below a comment
2. Enter reply content
3. Click "Post Reply"

#### @ Mentions

- Enter `@username` in comments
- Mentioned user receives a notification

### 14.5 Social Interactions

#### Liking

- Like a post: Click "Like" on post detail page
- Like a comment: Click "Like" below a comment

#### Bookmarking

1. Click "Bookmark" on post detail page
2. View in "Profile Center" → "My Bookmarks"

#### Following Users

1. Click a user avatar/username
2. Enter user profile page
3. Click the "Follow" button
4. Followed user earns 1 point

#### Unfollowing

1. Enter a followed user's profile page
2. Click the "Following" button

#### Private Messages

1. Click the "Messages" button at the top
2. Select a conversation or start a new private message
3. Enter message and send

#### Notifications

1. Click the "Notifications" button at the top
2. View four types of notifications:
   - 🔴 Like notifications (red)
   - 🔵 Comment notifications (blue)
   - 🟢 Follow notifications (green)
   - 🟡 Badge notifications (yellow)

### 14.6 User Profile Center

#### Accessing Profile Center

1. Click the user avatar at top right
2. Select "Profile Center"

#### Profile Information

**Viewing**:
- Avatar
- Nickname
- Bio
- Personal website link
- Statistics

**Editing**:
1. Click "Edit Profile"
2. Modify nickname, bio, website link
3. Click "Save"

#### My Achievements

**Check-in**:
1. Click "Check-in" on right sidebar of home page
2. Earn 10 points daily
3. Consecutive check-ins earn extra 5 points

**Badges**:
- View earned badges
- View requirements for unearned badges

#### My Posts

- View own posts
- Edit or delete posts

#### My Bookmarks

- View bookmarked posts
- Remove bookmarks

#### My Following

- View followed users
- Unfollow

#### My Followers

- View users following you
- Follow back or remove

---

## Chapter 15: FAQ and Troubleshooting

### 15.1 Admin Panel Issues

#### Q: Cannot login to admin panel?

**A**: Please check:
1. Username and password are correct
2. Admin panel URL is correct
3. Backend service is running normally
4. Browser console has error messages

#### Q: Configuration not taking effect after saving?

**A**: Try:
1. Refresh admin panel page
2. Clear browser cache
3. Check backend logs
4. Confirm configuration was actually saved successfully

#### Q: List loading slowly?

**A**: Possible causes:
1. Too much data, recommend paginated queries
2. Slow database queries, check indexes
3. Network issues
4. High backend service load

### 15.2 Main Site Issues

#### Q: Users cannot register?

**A**: Please check:
1. "Allow User Registration" is enabled in System Settings
2. Username doesn't already exist
3. Email format is correct
4. Password length meets requirements

#### Q: Users cannot post?

**A**: Please check:
1. "Allow Posting" is enabled in System Settings
2. User is not muted
3. User's reputation score is not too low
4. Not within rate limit period

#### Q: Image upload failing?

**A**: Please check:
1. Storage configuration is correct (Qiniu/Alibaba Cloud/Tencent Cloud)
2. File size is not over limit
3. File format is supported
4. Network connection is normal

### 15.3 Database Issues

#### Q: Where is the database file?

**A**: SQLite database file is typically in:
- `data` folder in project root
- File named `bbs.db` or similar

#### Q: How to backup database?

**A**: Directly copy database file:
1. Stop backend service
2. Copy `data/bbs.db` to backup location
3. Restart backend service

#### Q: Database corrupted?

**A**: Try:
1. Restore from backup
2. Use SQLite tools to check and repair
3. Reinitialize (will lose data)

---

## Chapter 16: Best Practices

### 16.1 Daily Operation Recommendations

#### Daily Checklist

- [ ] View pending reports
- [ ] View new users
- [ ] View hot posts
- [ ] Check system status
- [ ] View statistics

#### Weekly Checklist

- [ ] Analyze user growth trends
- [ ] Analyze content quality trends
- [ ] Check anti-spam system effectiveness
- [ ] Adjust configuration parameters
- [ ] Backup database

#### Monthly Checklist

- [ ] Comprehensive review of operation data
- [ ] Adjust community strategy
- [ ] Optimize board structure
- [ ] Update announcements
- [ ] Large-scale data backup

### 16.2 Community Management Strategies

#### Content Management

- **Promptly handle reports**: Process reports within 24 hours
- **Feature quality content**: Mark quality content as featured promptly
- **Pin important information**: Pin important announcements and events promptly
- **Regularly clean violations**: Regularly clean up violation content and users

#### User Management

- **Welcome new users**: Proactively welcome new users after registration
- **Cultivate core users**: Pay attention to active users, cultivate moderators
- **Mediate conflicts promptly**: Mediate user conflicts promptly
- **Fair and impartial**: Maintain fairness and impartiality when handling issues

#### Event Operation

- **Regularly hold events**: Hold online events weekly/monthly
- **Set reward mechanisms**: Set point/badge rewards for events
- **Promote**: Promote through announcements, pinning, etc.
- **Summarize and review**: Summarize experience after event ends

### 16.3 Security Recommendations

#### Account Security

- **Use strong passwords**: Use strong passwords for administrator accounts
- **Change passwords regularly**: Change password every 3-6 months
- **Don't share accounts**: Don't share administrator account with others
- **Enable 2FA**: Enable two-factor authentication if supported

#### System Security

- **Update regularly**: Update system and dependencies promptly
- **Backup data**: Regularly backup database
- **Restrict access**: Restrict IP access to admin panel
- **Monitor logs**: Regularly check system logs

#### Content Security

- **Moderation mechanism**: Establish content moderation mechanism
- **Anti-spam configuration**: Reasonably configure anti-spam parameters
- **Report handling**: Promptly handle user reports
- **User education**: Educate users to follow community rules

---

## Conclusion

Thank you for using the RainbowBBS community forum system!

We hope this operation guide helps you smoothly manage your community. If you have any questions or suggestions during use, welcome to:

- Submit an Issue to the GitHub repository
- Contact technical support
- Participate in community discussions

Wishing your community flourishes!

---

*Last updated: 2024*

