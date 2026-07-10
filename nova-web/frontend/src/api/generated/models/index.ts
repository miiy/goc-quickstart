/* tslint:disable */
/* eslint-disable */
/**
 * 
 * @export
 * @interface AuthUser
 */
export interface AuthUser {
    /**
     * 
     * @type {string}
     * @memberof AuthUser
     */
    username: string;
    /**
     * 
     * @type {number}
     * @memberof AuthUser
     */
    id: number;
}
/**
 * Read-only post category.
 * @export
 * @interface Category
 */
export interface Category {
    /**
     * 
     * @type {number}
     * @memberof Category
     */
    id: number;
    /**
     * 
     * @type {string}
     * @memberof Category
     */
    name: string;
    /**
     * 
     * @type {number}
     * @memberof Category
     */
    parentId: number;
    /**
     * 
     * @type {string}
     * @memberof Category
     */
    path: string;
}
/**
 * 
 * @export
 * @interface ChangePasswordRequest
 */
export interface ChangePasswordRequest {
    /**
     * 
     * @type {string}
     * @memberof ChangePasswordRequest
     */
    oldPassword: string;
    /**
     * 
     * @type {string}
     * @memberof ChangePasswordRequest
     */
    newPassword: string;
    /**
     * 
     * @type {string}
     * @memberof ChangePasswordRequest
     */
    newPasswordConfirmation: string;
}
/**
 * Post fields accepted when creating a post.
 * @export
 * @interface CreatePostInput
 */
export interface CreatePostInput {
    /**
     * 
     * @type {string}
     * @memberof CreatePostInput
     */
    title: string;
    /**
     * 
     * @type {string}
     * @memberof CreatePostInput
     */
    summary?: string;
    /**
     * 
     * @type {string}
     * @memberof CreatePostInput
     */
    content: string;
    /**
     * 
     * @type {PostStatus}
     * @memberof CreatePostInput
     */
    status?: PostStatus;
    /**
     * 
     * @type {Array<string>}
     * @memberof CreatePostInput
     */
    tags?: Array<string>;
    /**
     * 
     * @type {number}
     * @memberof CreatePostInput
     */
    categoryId?: number;
    /**
     * Cover image object key or /uploads URL returned by file upload.
     * @type {string}
     * @memberof CreatePostInput
     */
    coverUrl?: string;
    /**
     * 
     * @type {string}
     * @memberof CreatePostInput
     */
    publishedAt?: string | null;
}


/**
 * 
 * @export
 * @interface CreatePostRequest
 */
export interface CreatePostRequest {
    /**
     * 
     * @type {CreatePostInput}
     * @memberof CreatePostRequest
     */
    post: CreatePostInput;
}
/**
 * 
 * @export
 * @interface CreatePostResponse
 */
export interface CreatePostResponse {
    /**
     * 
     * @type {Post}
     * @memberof CreatePostResponse
     */
    post: Post;
}
/**
 * 
 * @export
 * @interface EmailCheckRequest
 */
export interface EmailCheckRequest {
    /**
     * 
     * @type {string}
     * @memberof EmailCheckRequest
     */
    value: string;
}
/**
 * 
 * @export
 * @interface EmailCheckResponse
 */
export interface EmailCheckResponse {
    /**
     * 
     * @type {boolean}
     * @memberof EmailCheckResponse
     */
    exist: boolean;
}
/**
 * 
 * @export
 * @interface ErrorResponse
 */
export interface ErrorResponse {
    /**
     * 
     * @type {ErrorStatus}
     * @memberof ErrorResponse
     */
    error: ErrorStatus;
}
/**
 * 
 * @export
 * @interface ErrorStatus
 */
export interface ErrorStatus {
    /**
     * 
     * @type {number}
     * @memberof ErrorStatus
     */
    code: number;
    /**
     * 
     * @type {string}
     * @memberof ErrorStatus
     */
    message: string;
}
/**
 * Business scene for a stored file.
 * @export
 * @enum {string}
 */
export enum FileScene {
    FILE_SCENE_UNSPECIFIED = 'unspecified',
    FILE_SCENE_AVATAR = 'avatar',
    FILE_SCENE_POST_COVER = 'post_cover',
    FILE_SCENE_POST_CONTENT = 'post_content'
}

/**
 * Lifecycle state of a stored file.
 * @export
 * @enum {string}
 */
export enum FileStatus {
    FILE_STATUS_UNSPECIFIED = 'unspecified',
    FILE_STATUS_ACTIVE = 'active',
    FILE_STATUS_DELETED = 'deleted'
}

/**
 * 
 * @export
 * @interface GetPostResponse
 */
export interface GetPostResponse {
    /**
     * 
     * @type {Post}
     * @memberof GetPostResponse
     */
    post: Post;
}
/**
 * 
 * @export
 * @interface GetProfileResponse
 */
export interface GetProfileResponse {
    /**
     * 
     * @type {User}
     * @memberof GetProfileResponse
     */
    user: User;
}
/**
 * 
 * @export
 * @interface GetUserResponse
 */
export interface GetUserResponse {
    /**
     * 
     * @type {PublicUser}
     * @memberof GetUserResponse
     */
    user: PublicUser;
}
/**
 * 
 * @export
 * @interface ListCategoriesResponse
 */
export interface ListCategoriesResponse {
    /**
     * 
     * @type {Array<Category>}
     * @memberof ListCategoriesResponse
     */
    categories: Array<Category>;
}
/**
 * 
 * @export
 * @interface ListPostsResponse
 */
export interface ListPostsResponse {
    /**
     * 
     * @type {number}
     * @memberof ListPostsResponse
     */
    total: number;
    /**
     * 
     * @type {number}
     * @memberof ListPostsResponse
     */
    totalPages: number;
    /**
     * 
     * @type {number}
     * @memberof ListPostsResponse
     */
    pageSize: number;
    /**
     * 
     * @type {number}
     * @memberof ListPostsResponse
     */
    currentPage: number;
    /**
     * 
     * @type {Array<Post>}
     * @memberof ListPostsResponse
     */
    posts: Array<Post>;
}
/**
 * 
 * @export
 * @interface ListUsersResponse
 */
export interface ListUsersResponse {
    /**
     * 
     * @type {number}
     * @memberof ListUsersResponse
     */
    total: number;
    /**
     * 
     * @type {number}
     * @memberof ListUsersResponse
     */
    totalPages: number;
    /**
     * 
     * @type {number}
     * @memberof ListUsersResponse
     */
    pageSize: number;
    /**
     * 
     * @type {number}
     * @memberof ListUsersResponse
     */
    currentPage: number;
    /**
     * 
     * @type {Array<User>}
     * @memberof ListUsersResponse
     */
    users: Array<User>;
}
/**
 * 
 * @export
 * @interface LoginRequest
 */
export interface LoginRequest {
    /**
     * 
     * @type {string}
     * @memberof LoginRequest
     */
    username: string;
    /**
     * 
     * @type {string}
     * @memberof LoginRequest
     */
    password: string;
}
/**
 * Either access_token or refresh_token is required. If access_token is omitted, the gateway uses the bearer token.
 * @export
 * @interface LogoutRequest
 */
export interface LogoutRequest {
    /**
     * 
     * @type {string}
     * @memberof LogoutRequest
     */
    accessToken?: string;
    /**
     * 
     * @type {string}
     * @memberof LogoutRequest
     */
    refreshToken?: string;
}
/**
 * File stored by the file service.
 * @export
 * @interface ModelFile
 */
export interface ModelFile {
    /**
     * 
     * @type {number}
     * @memberof ModelFile
     */
    id: number;
    /**
     * 
     * @type {number}
     * @memberof ModelFile
     */
    ownerId: number;
    /**
     * 
     * @type {string}
     * @memberof ModelFile
     */
    ownerType: string;
    /**
     * 
     * @type {FileScene}
     * @memberof ModelFile
     */
    scene: FileScene;
    /**
     * Stable relative storage key, for example avatars/2026/06/a.png.
     * @type {string}
     * @memberof ModelFile
     */
    objectKey: string;
    /**
     * 
     * @type {string}
     * @memberof ModelFile
     */
    url: string;
    /**
     * 
     * @type {string}
     * @memberof ModelFile
     */
    mimeType: string;
    /**
     * 
     * @type {number}
     * @memberof ModelFile
     */
    size: number;
    /**
     * 
     * @type {string}
     * @memberof ModelFile
     */
    checksum: string;
    /**
     * 
     * @type {FileStatus}
     * @memberof ModelFile
     */
    status: FileStatus;
    /**
     * 
     * @type {number}
     * @memberof ModelFile
     */
    createdBy: number;
    /**
     * 
     * @type {string}
     * @memberof ModelFile
     */
    createdAt: string;
    /**
     * 
     * @type {string}
     * @memberof ModelFile
     */
    updatedAt: string;
    /**
     * 
     * @type {string}
     * @memberof ModelFile
     */
    deletedAt?: string | null;
}


/**
 * 
 * @export
 * @interface MpLoginRequest
 */
export interface MpLoginRequest {
    /**
     * 
     * @type {string}
     * @memberof MpLoginRequest
     */
    code: string;
}
/**
 * 
 * @export
 * @interface PhoneAuthRequest
 */
export interface PhoneAuthRequest {
    /**
     * 
     * @type {string}
     * @memberof PhoneAuthRequest
     */
    phone: string;
    /**
     * 
     * @type {string}
     * @memberof PhoneAuthRequest
     */
    code: string;
}
/**
 * 
 * @export
 * @interface PhoneCheckRequest
 */
export interface PhoneCheckRequest {
    /**
     * 
     * @type {string}
     * @memberof PhoneCheckRequest
     */
    value: string;
}
/**
 * 
 * @export
 * @interface PhoneCheckResponse
 */
export interface PhoneCheckResponse {
    /**
     * 
     * @type {boolean}
     * @memberof PhoneCheckResponse
     */
    exist: boolean;
}
/**
 * Article published by a user.
 * @export
 * @interface Post
 */
export interface Post {
    /**
     * Public post id used in URLs.
     * @type {string}
     * @memberof Post
     */
    id: string;
    /**
     * 
     * @type {number}
     * @memberof Post
     */
    userId: number;
    /**
     * 
     * @type {PostUser}
     * @memberof Post
     */
    user: PostUser;
    /**
     * 
     * @type {string}
     * @memberof Post
     */
    title: string;
    /**
     * 
     * @type {string}
     * @memberof Post
     */
    summary?: string;
    /**
     * 
     * @type {string}
     * @memberof Post
     */
    content: string;
    /**
     * 
     * @type {PostStatus}
     * @memberof Post
     */
    status: PostStatus;
    /**
     * 
     * @type {Array<string>}
     * @memberof Post
     */
    tags: Array<string>;
    /**
     * 
     * @type {number}
     * @memberof Post
     */
    categoryId: number;
    /**
     * 
     * @type {string}
     * @memberof Post
     */
    publishedAt?: string | null;
    /**
     * 
     * @type {string}
     * @memberof Post
     */
    createdAt: string;
    /**
     * 
     * @type {string}
     * @memberof Post
     */
    updatedAt: string;
    /**
     * 
     * @type {string}
     * @memberof Post
     */
    deletedAt?: string | null;
    /**
     * Browser-readable cover image URL, for example /uploads/post-covers/2026/07/a.png.
     * @type {string}
     * @memberof Post
     */
    coverUrl: string;
    /**
     * Whether the authenticated caller can manage this post.
     * @type {boolean}
     * @memberof Post
     */
    canManage?: boolean;
}


/**
 * Publication state of a post.
 * @export
 * @enum {string}
 */
export enum PostStatus {
    POST_STATUS_UNSPECIFIED = 'unspecified',
    POST_STATUS_DRAFT = 'draft',
    POST_STATUS_PUBLISHED = 'published',
    POST_STATUS_PENDING_REVIEW = 'pending_review'
}

/**
 * User display fields populated by nova-gateway from nova-user.
 * @export
 * @interface PostUser
 */
export interface PostUser {
    /**
     * 
     * @type {string}
     * @memberof PostUser
     */
    username: string;
    /**
     * 
     * @type {string}
     * @memberof PostUser
     */
    nickname: string;
    /**
     * Browser-readable avatar URL, for example /uploads/avatars/2026/07/a.png.
     * @type {string}
     * @memberof PostUser
     */
    avatar: string;
}
/**
 * Public user profile fields safe for anonymous reads.
 * @export
 * @interface PublicUser
 */
export interface PublicUser {
    /**
     * 
     * @type {number}
     * @memberof PublicUser
     */
    id: number;
    /**
     * 
     * @type {string}
     * @memberof PublicUser
     */
    username: string;
    /**
     * 
     * @type {string}
     * @memberof PublicUser
     */
    nickname: string;
    /**
     * Browser-readable avatar URL, for example /uploads/avatars/2026/07/a.png.
     * @type {string}
     * @memberof PublicUser
     */
    avatar: string;
    /**
     * 
     * @type {string}
     * @memberof PublicUser
     */
    createdAt: string;
    /**
     * 
     * @type {string}
     * @memberof PublicUser
     */
    updatedAt: string;
}
/**
 * 
 * @export
 * @interface RefreshTokenRequest
 */
export interface RefreshTokenRequest {
    /**
     * 
     * @type {string}
     * @memberof RefreshTokenRequest
     */
    refreshToken: string;
}
/**
 * 
 * @export
 * @interface RegisterRequest
 */
export interface RegisterRequest {
    /**
     * 
     * @type {string}
     * @memberof RegisterRequest
     */
    email: string;
    /**
     * 
     * @type {string}
     * @memberof RegisterRequest
     */
    username: string;
    /**
     * 
     * @type {string}
     * @memberof RegisterRequest
     */
    password: string;
    /**
     * 
     * @type {string}
     * @memberof RegisterRequest
     */
    passwordConfirmation: string;
}
/**
 * 
 * @export
 * @interface RegisterResponse
 */
export interface RegisterResponse {
    /**
     * 
     * @type {AuthUser}
     * @memberof RegisterResponse
     */
    user: AuthUser;
}
/**
 * 
 * @export
 * @interface SendSmsCodeRequest
 */
export interface SendSmsCodeRequest {
    /**
     * 
     * @type {string}
     * @memberof SendSmsCodeRequest
     */
    phone: string;
}
/**
 * 
 * @export
 * @interface TokenResponse
 */
export interface TokenResponse {
    /**
     * 
     * @type {string}
     * @memberof TokenResponse
     */
    tokenType: string;
    /**
     * 
     * @type {string}
     * @memberof TokenResponse
     */
    accessToken: string;
    /**
     * 
     * @type {string}
     * @memberof TokenResponse
     */
    expiresAt: string;
    /**
     * 
     * @type {AuthUser}
     * @memberof TokenResponse
     */
    user: AuthUser;
    /**
     * 
     * @type {string}
     * @memberof TokenResponse
     */
    refreshToken: string;
    /**
     * 
     * @type {string}
     * @memberof TokenResponse
     */
    refreshExpiresAt: string;
}
/**
 * Mutable post fields accepted when updating a post.
 * @export
 * @interface UpdatePostInput
 */
export interface UpdatePostInput {
    /**
     * 
     * @type {string}
     * @memberof UpdatePostInput
     */
    title?: string;
    /**
     * 
     * @type {string}
     * @memberof UpdatePostInput
     */
    summary?: string;
    /**
     * 
     * @type {string}
     * @memberof UpdatePostInput
     */
    content?: string;
    /**
     * 
     * @type {PostStatus}
     * @memberof UpdatePostInput
     */
    status?: PostStatus;
    /**
     * 
     * @type {Array<string>}
     * @memberof UpdatePostInput
     */
    tags?: Array<string>;
    /**
     * 
     * @type {number}
     * @memberof UpdatePostInput
     */
    categoryId?: number;
    /**
     * Cover image object key or /uploads URL returned by file upload.
     * @type {string}
     * @memberof UpdatePostInput
     */
    coverUrl?: string;
    /**
     * 
     * @type {string}
     * @memberof UpdatePostInput
     */
    publishedAt?: string | null;
}


/**
 * 
 * @export
 * @interface UpdatePostRequest
 */
export interface UpdatePostRequest {
    /**
     * 
     * @type {UpdatePostInput}
     * @memberof UpdatePostRequest
     */
    post: UpdatePostInput;
    /**
     * Mutable post fields to update.
     * @type {Array<UpdatePostRequestUpdateFieldsEnum>}
     * @memberof UpdatePostRequest
     */
    updateFields?: Array<UpdatePostRequestUpdateFieldsEnum>;
}

/**
* @export
* @enum {string}
*/
export enum UpdatePostRequestUpdateFieldsEnum {
    Title = 'title',
    Summary = 'summary',
    Content = 'content',
    CoverUrl = 'cover_url',
    PublishedAt = 'published_at',
    Status = 'status',
    Tags = 'tags',
    CategoryId = 'category_id'
}

/**
 * 
 * @export
 * @interface UpdatePostResponse
 */
export interface UpdatePostResponse {
    /**
     * 
     * @type {Post}
     * @memberof UpdatePostResponse
     */
    post: Post;
}
/**
 * 
 * @export
 * @interface UpdateProfileRequest
 */
export interface UpdateProfileRequest {
    /**
     * 
     * @type {UserInput}
     * @memberof UpdateProfileRequest
     */
    user: UserInput;
    /**
     * Mutable user fields to update.
     * @type {Array<UpdateProfileRequestUpdateFieldsEnum>}
     * @memberof UpdateProfileRequest
     */
    updateFields?: Array<UpdateProfileRequestUpdateFieldsEnum>;
}

/**
* @export
* @enum {string}
*/
export enum UpdateProfileRequestUpdateFieldsEnum {
    Nickname = 'nickname',
    Avatar = 'avatar',
    Email = 'email',
    Phone = 'phone',
    Status = 'status'
}

/**
 * 
 * @export
 * @interface UpdateProfileResponse
 */
export interface UpdateProfileResponse {
    /**
     * 
     * @type {User}
     * @memberof UpdateProfileResponse
     */
    user: User;
}
/**
 * 
 * @export
 * @interface UploadFileResponse
 */
export interface UploadFileResponse {
    /**
     * 
     * @type {any}
     * @memberof UploadFileResponse
     */
    file: any;
}
/**
 * Public profile and account state for a user.
 * @export
 * @interface User
 */
export interface User {
    /**
     * 
     * @type {number}
     * @memberof User
     */
    id: number;
    /**
     * 
     * @type {string}
     * @memberof User
     */
    username: string;
    /**
     * 
     * @type {string}
     * @memberof User
     */
    nickname: string;
    /**
     * Browser-readable avatar URL, for example /uploads/avatars/2026/07/a.png.
     * @type {string}
     * @memberof User
     */
    avatar: string;
    /**
     * 
     * @type {string}
     * @memberof User
     */
    email: string;
    /**
     * 
     * @type {string}
     * @memberof User
     */
    emailVerifiedTime?: string | null;
    /**
     * 
     * @type {string}
     * @memberof User
     */
    phone: string;
    /**
     * 
     * @type {UserStatus}
     * @memberof User
     */
    status: UserStatus;
    /**
     * 
     * @type {string}
     * @memberof User
     */
    createdAt: string;
    /**
     * 
     * @type {string}
     * @memberof User
     */
    updatedAt: string;
    /**
     * 
     * @type {string}
     * @memberof User
     */
    deletedAt?: string | null;
}


/**
 * Mutable user profile fields accepted by the gateway.
 * @export
 * @interface UserInput
 */
export interface UserInput {
    /**
     * 
     * @type {string}
     * @memberof UserInput
     */
    nickname?: string;
    /**
     * Avatar object key or /uploads URL returned by avatar upload.
     * @type {string}
     * @memberof UserInput
     */
    avatar?: string;
    /**
     * 
     * @type {string}
     * @memberof UserInput
     */
    email?: string;
    /**
     * 
     * @type {string}
     * @memberof UserInput
     */
    phone?: string;
    /**
     * 
     * @type {UserStatus}
     * @memberof UserInput
     */
    status?: UserStatus;
}


/**
 * Lifecycle state of a user account.
 * @export
 * @enum {string}
 */
export enum UserStatus {
    USER_STATUS_UNSPECIFIED = 'unspecified',
    USER_STATUS_ACTIVE = 'active',
    USER_STATUS_DISABLED = 'disabled'
}

/**
 * 
 * @export
 * @interface UsernameCheckRequest
 */
export interface UsernameCheckRequest {
    /**
     * 
     * @type {string}
     * @memberof UsernameCheckRequest
     */
    value: string;
}
/**
 * 
 * @export
 * @interface UsernameCheckResponse
 */
export interface UsernameCheckResponse {
    /**
     * 
     * @type {boolean}
     * @memberof UsernameCheckResponse
     */
    exist: boolean;
}
