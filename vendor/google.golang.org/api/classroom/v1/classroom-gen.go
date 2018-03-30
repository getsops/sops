// Package classroom provides access to the Google Classroom API.
//
// See https://developers.google.com/classroom/
//
// Usage example:
//
//   import "google.golang.org/api/classroom/v1"
//   ...
//   classroomService, err := classroom.New(oauthHttpClient)
package classroom // import "google.golang.org/api/classroom/v1"

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	context "golang.org/x/net/context"
	ctxhttp "golang.org/x/net/context/ctxhttp"
	gensupport "google.golang.org/api/gensupport"
	googleapi "google.golang.org/api/googleapi"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Always reference these packages, just in case the auto-generated code
// below doesn't.
var _ = bytes.NewBuffer
var _ = strconv.Itoa
var _ = fmt.Sprintf
var _ = json.NewDecoder
var _ = io.Copy
var _ = url.Parse
var _ = gensupport.MarshalJSON
var _ = googleapi.Version
var _ = errors.New
var _ = strings.Replace
var _ = context.Canceled
var _ = ctxhttp.Do

const apiId = "classroom:v1"
const apiName = "classroom"
const apiVersion = "v1"
const basePath = "https://classroom.googleapis.com/"

// OAuth2 scopes used by this API.
const (
	// View and manage announcements in Google Classroom
	ClassroomAnnouncementsScope = "https://www.googleapis.com/auth/classroom.announcements"

	// View announcements in Google Classroom
	ClassroomAnnouncementsReadonlyScope = "https://www.googleapis.com/auth/classroom.announcements.readonly"

	// Manage your Google Classroom classes
	ClassroomCoursesScope = "https://www.googleapis.com/auth/classroom.courses"

	// View your Google Classroom classes
	ClassroomCoursesReadonlyScope = "https://www.googleapis.com/auth/classroom.courses.readonly"

	// Manage your course work and view your grades in Google Classroom
	ClassroomCourseworkMeScope = "https://www.googleapis.com/auth/classroom.coursework.me"

	// View your course work and grades in Google Classroom
	ClassroomCourseworkMeReadonlyScope = "https://www.googleapis.com/auth/classroom.coursework.me.readonly"

	// Manage course work and grades for students in the Google Classroom
	// classes you teach and view the course work and grades for classes you
	// administer
	ClassroomCourseworkStudentsScope = "https://www.googleapis.com/auth/classroom.coursework.students"

	// View course work and grades for students in the Google Classroom
	// classes you teach or administer
	ClassroomCourseworkStudentsReadonlyScope = "https://www.googleapis.com/auth/classroom.coursework.students.readonly"

	// View your Google Classroom guardians
	ClassroomGuardianlinksMeReadonlyScope = "https://www.googleapis.com/auth/classroom.guardianlinks.me.readonly"

	// View and manage guardians for students in your Google Classroom
	// classes
	ClassroomGuardianlinksStudentsScope = "https://www.googleapis.com/auth/classroom.guardianlinks.students"

	// View guardians for students in your Google Classroom classes
	ClassroomGuardianlinksStudentsReadonlyScope = "https://www.googleapis.com/auth/classroom.guardianlinks.students.readonly"

	// View the email addresses of people in your classes
	ClassroomProfileEmailsScope = "https://www.googleapis.com/auth/classroom.profile.emails"

	// View the profile photos of people in your classes
	ClassroomProfilePhotosScope = "https://www.googleapis.com/auth/classroom.profile.photos"

	// Manage your Google Classroom class rosters
	ClassroomRostersScope = "https://www.googleapis.com/auth/classroom.rosters"

	// View your Google Classroom class rosters
	ClassroomRostersReadonlyScope = "https://www.googleapis.com/auth/classroom.rosters.readonly"

	// View your course work and grades in Google Classroom
	ClassroomStudentSubmissionsMeReadonlyScope = "https://www.googleapis.com/auth/classroom.student-submissions.me.readonly"

	// View course work and grades for students in the Google Classroom
	// classes you teach or administer
	ClassroomStudentSubmissionsStudentsReadonlyScope = "https://www.googleapis.com/auth/classroom.student-submissions.students.readonly"
)

func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	s := &Service{client: client, BasePath: basePath}
	s.Courses = NewCoursesService(s)
	s.Invitations = NewInvitationsService(s)
	s.Registrations = NewRegistrationsService(s)
	s.UserProfiles = NewUserProfilesService(s)
	return s, nil
}

type Service struct {
	client    *http.Client
	BasePath  string // API endpoint base URL
	UserAgent string // optional additional User-Agent fragment

	Courses *CoursesService

	Invitations *InvitationsService

	Registrations *RegistrationsService

	UserProfiles *UserProfilesService
}

func (s *Service) userAgent() string {
	if s.UserAgent == "" {
		return googleapi.UserAgent
	}
	return googleapi.UserAgent + " " + s.UserAgent
}

func NewCoursesService(s *Service) *CoursesService {
	rs := &CoursesService{s: s}
	rs.Aliases = NewCoursesAliasesService(s)
	rs.Announcements = NewCoursesAnnouncementsService(s)
	rs.CourseWork = NewCoursesCourseWorkService(s)
	rs.Students = NewCoursesStudentsService(s)
	rs.Teachers = NewCoursesTeachersService(s)
	return rs
}

type CoursesService struct {
	s *Service

	Aliases *CoursesAliasesService

	Announcements *CoursesAnnouncementsService

	CourseWork *CoursesCourseWorkService

	Students *CoursesStudentsService

	Teachers *CoursesTeachersService
}

func NewCoursesAliasesService(s *Service) *CoursesAliasesService {
	rs := &CoursesAliasesService{s: s}
	return rs
}

type CoursesAliasesService struct {
	s *Service
}

func NewCoursesAnnouncementsService(s *Service) *CoursesAnnouncementsService {
	rs := &CoursesAnnouncementsService{s: s}
	return rs
}

type CoursesAnnouncementsService struct {
	s *Service
}

func NewCoursesCourseWorkService(s *Service) *CoursesCourseWorkService {
	rs := &CoursesCourseWorkService{s: s}
	rs.StudentSubmissions = NewCoursesCourseWorkStudentSubmissionsService(s)
	return rs
}

type CoursesCourseWorkService struct {
	s *Service

	StudentSubmissions *CoursesCourseWorkStudentSubmissionsService
}

func NewCoursesCourseWorkStudentSubmissionsService(s *Service) *CoursesCourseWorkStudentSubmissionsService {
	rs := &CoursesCourseWorkStudentSubmissionsService{s: s}
	return rs
}

type CoursesCourseWorkStudentSubmissionsService struct {
	s *Service
}

func NewCoursesStudentsService(s *Service) *CoursesStudentsService {
	rs := &CoursesStudentsService{s: s}
	return rs
}

type CoursesStudentsService struct {
	s *Service
}

func NewCoursesTeachersService(s *Service) *CoursesTeachersService {
	rs := &CoursesTeachersService{s: s}
	return rs
}

type CoursesTeachersService struct {
	s *Service
}

func NewInvitationsService(s *Service) *InvitationsService {
	rs := &InvitationsService{s: s}
	return rs
}

type InvitationsService struct {
	s *Service
}

func NewRegistrationsService(s *Service) *RegistrationsService {
	rs := &RegistrationsService{s: s}
	return rs
}

type RegistrationsService struct {
	s *Service
}

func NewUserProfilesService(s *Service) *UserProfilesService {
	rs := &UserProfilesService{s: s}
	rs.GuardianInvitations = NewUserProfilesGuardianInvitationsService(s)
	rs.Guardians = NewUserProfilesGuardiansService(s)
	return rs
}

type UserProfilesService struct {
	s *Service

	GuardianInvitations *UserProfilesGuardianInvitationsService

	Guardians *UserProfilesGuardiansService
}

func NewUserProfilesGuardianInvitationsService(s *Service) *UserProfilesGuardianInvitationsService {
	rs := &UserProfilesGuardianInvitationsService{s: s}
	return rs
}

type UserProfilesGuardianInvitationsService struct {
	s *Service
}

func NewUserProfilesGuardiansService(s *Service) *UserProfilesGuardiansService {
	rs := &UserProfilesGuardiansService{s: s}
	return rs
}

type UserProfilesGuardiansService struct {
	s *Service
}

// Announcement: Announcement created by a teacher for students of the
// course
type Announcement struct {
	// AlternateLink: Absolute link to this announcement in the Classroom
	// web UI.
	// This is only populated if `state` is `PUBLISHED`.
	//
	// Read-only.
	AlternateLink string `json:"alternateLink,omitempty"`

	// AssigneeMode: Assignee mode of the announcement.
	// If unspecified, the default value is `ALL_STUDENTS`.
	//
	// Possible values:
	//   "ASSIGNEE_MODE_UNSPECIFIED" - No mode specified. This is never
	// returned.
	//   "ALL_STUDENTS" - All students can see the item.
	// This is the default state.
	//   "INDIVIDUAL_STUDENTS" - A subset of the students can see the item.
	AssigneeMode string `json:"assigneeMode,omitempty"`

	// CourseId: Identifier of the course.
	//
	// Read-only.
	CourseId string `json:"courseId,omitempty"`

	// CreationTime: Timestamp when this announcement was
	// created.
	//
	// Read-only.
	CreationTime string `json:"creationTime,omitempty"`

	// CreatorUserId: Identifier for the user that created the
	// announcement.
	//
	// Read-only.
	CreatorUserId string `json:"creatorUserId,omitempty"`

	// Id: Classroom-assigned identifier of this announcement, unique per
	// course.
	//
	// Read-only.
	Id string `json:"id,omitempty"`

	// IndividualStudentsOptions: Identifiers of students with access to the
	// announcement.
	// This field is set only if `assigneeMode` is `INDIVIDUAL_STUDENTS`.
	// If the `assigneeMode` is `INDIVIDUAL_STUDENTS`, then only students
	// specified in this
	// field will be able to see the announcement.
	IndividualStudentsOptions *IndividualStudentsOptions `json:"individualStudentsOptions,omitempty"`

	// Materials: Additional materials.
	//
	// Announcements must have no more than 20 material items.
	Materials []*Material `json:"materials,omitempty"`

	// ScheduledTime: Optional timestamp when this announcement is scheduled
	// to be published.
	ScheduledTime string `json:"scheduledTime,omitempty"`

	// State: Status of this announcement.
	// If unspecified, the default state is `DRAFT`.
	//
	// Possible values:
	//   "ANNOUNCEMENT_STATE_UNSPECIFIED" - No state specified. This is
	// never returned.
	//   "PUBLISHED" - Status for announcement that has been published.
	// This is the default state.
	//   "DRAFT" - Status for an announcement that is not yet
	// published.
	// Announcement in this state is visible only to course teachers and
	// domain
	// administrators.
	//   "DELETED" - Status for announcement that was published but is now
	// deleted.
	// Announcement in this state is visible only to course teachers and
	// domain
	// administrators.
	// Announcement in this state is deleted after some time.
	State string `json:"state,omitempty"`

	// Text: Description of this announcement.
	// The text must be a valid UTF-8 string containing no more
	// than 30,000 characters.
	Text string `json:"text,omitempty"`

	// UpdateTime: Timestamp of the most recent change to this
	// announcement.
	//
	// Read-only.
	UpdateTime string `json:"updateTime,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "AlternateLink") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AlternateLink") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Announcement) MarshalJSON() ([]byte, error) {
	type noMethod Announcement
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Assignment: Additional details for assignments.
type Assignment struct {
	// StudentWorkFolder: Drive folder where attachments from student
	// submissions are placed.
	// This is only populated for course teachers and administrators.
	StudentWorkFolder *DriveFolder `json:"studentWorkFolder,omitempty"`

	// ForceSendFields is a list of field names (e.g. "StudentWorkFolder")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "StudentWorkFolder") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *Assignment) MarshalJSON() ([]byte, error) {
	type noMethod Assignment
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// AssignmentSubmission: Student work for an assignment.
type AssignmentSubmission struct {
	// Attachments: Attachments added by the student.
	// Drive files that correspond to materials with a share mode
	// of
	// STUDENT_COPY may not exist yet if the student has not accessed
	// the
	// assignment in Classroom.
	//
	// Some attachment metadata is only populated if the requesting user
	// has
	// permission to access it. Identifier and alternate_link fields are
	// always
	// available, but others (e.g. title) may not be.
	Attachments []*Attachment `json:"attachments,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Attachments") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Attachments") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *AssignmentSubmission) MarshalJSON() ([]byte, error) {
	type noMethod AssignmentSubmission
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Attachment: Attachment added to student assignment work.
//
// When creating attachments, setting the `form` field is not supported.
type Attachment struct {
	// DriveFile: Google Drive file attachment.
	DriveFile *DriveFile `json:"driveFile,omitempty"`

	// Form: Google Forms attachment.
	Form *Form `json:"form,omitempty"`

	// Link: Link attachment.
	Link *Link `json:"link,omitempty"`

	// YouTubeVideo: Youtube video attachment.
	YouTubeVideo *YouTubeVideo `json:"youTubeVideo,omitempty"`

	// ForceSendFields is a list of field names (e.g. "DriveFile") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "DriveFile") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Attachment) MarshalJSON() ([]byte, error) {
	type noMethod Attachment
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CloudPubsubTopic: A reference to a Cloud Pub/Sub topic.
//
// To register for notifications, the owner of the topic must
// grant
// `classroom-notifications@system.gserviceaccount.com` the
//  `projects.topics.publish` permission.
type CloudPubsubTopic struct {
	// TopicName: The `name` field of a Cloud
	// Pub/Sub
	// [Topic](https://cloud.google.com/pubsub/docs/reference/rest/v1
	// /projects.topics#Topic).
	TopicName string `json:"topicName,omitempty"`

	// ForceSendFields is a list of field names (e.g. "TopicName") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "TopicName") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *CloudPubsubTopic) MarshalJSON() ([]byte, error) {
	type noMethod CloudPubsubTopic
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Course: A Course in Classroom.
type Course struct {
	// AlternateLink: Absolute link to this course in the Classroom web
	// UI.
	//
	// Read-only.
	AlternateLink string `json:"alternateLink,omitempty"`

	// CalendarId: The Calendar ID for a calendar that all course members
	// can see, to which
	// Classroom adds events for course work and announcements in the
	// course.
	//
	// Read-only.
	CalendarId string `json:"calendarId,omitempty"`

	// CourseGroupEmail: The email address of a Google group containing all
	// members of the course.
	// This group does not accept email and can only be used for
	// permissions.
	//
	// Read-only.
	CourseGroupEmail string `json:"courseGroupEmail,omitempty"`

	// CourseMaterialSets: Sets of materials that appear on the "about" page
	// of this course.
	//
	// Read-only.
	CourseMaterialSets []*CourseMaterialSet `json:"courseMaterialSets,omitempty"`

	// CourseState: State of the course.
	// If unspecified, the default state is `PROVISIONED`.
	//
	// Possible values:
	//   "COURSE_STATE_UNSPECIFIED" - No course state. No returned Course
	// message will use this value.
	//   "ACTIVE" - The course is active.
	//   "ARCHIVED" - The course has been archived. You cannot modify it
	// except to change it
	// to a different state.
	//   "PROVISIONED" - The course has been created, but not yet activated.
	// It is accessible by
	// the primary teacher and domain administrators, who may modify it
	// or
	// change it to the `ACTIVE` or `DECLINED` states.
	// A course may only be changed to `PROVISIONED` if it is in the
	// `DECLINED`
	// state.
	//   "DECLINED" - The course has been created, but declined. It is
	// accessible by the
	// course owner and domain administrators, though it will not
	// be
	// displayed in the web UI. You cannot modify the course except to
	// change it
	// to the `PROVISIONED` state.
	// A course may only be changed to `DECLINED` if it is in the
	// `PROVISIONED`
	// state.
	//   "SUSPENDED" - The course has been suspended. You cannot modify the
	// course, and only the
	// user identified by the `owner_id` can view the course.
	// A course may be placed in this state if it potentially violates
	// the
	// Terms of Service.
	CourseState string `json:"courseState,omitempty"`

	// CreationTime: Creation time of the course.
	// Specifying this field in a course update mask results in an
	// error.
	//
	// Read-only.
	CreationTime string `json:"creationTime,omitempty"`

	// Description: Optional description.
	// For example, "We'll be learning about the structure of
	// living
	// creatures from a combination of textbooks, guest lectures, and lab
	// work.
	// Expect to be excited!"
	// If set, this field must be a valid UTF-8 string and no longer than
	// 30,000
	// characters.
	Description string `json:"description,omitempty"`

	// DescriptionHeading: Optional heading for the description.
	// For example, "Welcome to 10th Grade Biology."
	// If set, this field must be a valid UTF-8 string and no longer than
	// 3600
	// characters.
	DescriptionHeading string `json:"descriptionHeading,omitempty"`

	// EnrollmentCode: Enrollment code to use when joining this
	// course.
	// Specifying this field in a course update mask results in an
	// error.
	//
	// Read-only.
	EnrollmentCode string `json:"enrollmentCode,omitempty"`

	// GuardiansEnabled: Whether or not guardian notifications are enabled
	// for this course.
	//
	// Read-only.
	GuardiansEnabled bool `json:"guardiansEnabled,omitempty"`

	// Id: Identifier for this course assigned by Classroom.
	//
	// When
	// creating a course,
	// you may optionally set this identifier to an
	// alias string in the
	// request to create a corresponding alias. The `id` is still assigned
	// by
	// Classroom and cannot be updated after the course is
	// created.
	//
	// Specifying this field in a course update mask results in an error.
	Id string `json:"id,omitempty"`

	// Name: Name of the course.
	// For example, "10th Grade Biology".
	// The name is required. It must be between 1 and 750 characters and a
	// valid
	// UTF-8 string.
	Name string `json:"name,omitempty"`

	// OwnerId: The identifier of the owner of a course.
	//
	// When specified as a parameter of a
	// create course request, this
	// field is required.
	// The identifier can be one of the following:
	//
	// * the numeric identifier for the user
	// * the email address of the user
	// * the string literal "me", indicating the requesting user
	//
	// This must be set in a create request. Admins can also specify this
	// field
	// in a patch course request to
	// transfer ownership. In other contexts, it is read-only.
	OwnerId string `json:"ownerId,omitempty"`

	// Room: Optional room location.
	// For example, "301".
	// If set, this field must be a valid UTF-8 string and no longer than
	// 650
	// characters.
	Room string `json:"room,omitempty"`

	// Section: Section of the course.
	// For example, "Period 2".
	// If set, this field must be a valid UTF-8 string and no longer than
	// 2800
	// characters.
	Section string `json:"section,omitempty"`

	// TeacherFolder: Information about a Drive Folder that is shared with
	// all teachers of the
	// course.
	//
	// This field will only be set for teachers of the course and domain
	// administrators.
	//
	// Read-only.
	TeacherFolder *DriveFolder `json:"teacherFolder,omitempty"`

	// TeacherGroupEmail: The email address of a Google group containing all
	// teachers of the course.
	// This group does not accept email and can only be used for
	// permissions.
	//
	// Read-only.
	TeacherGroupEmail string `json:"teacherGroupEmail,omitempty"`

	// UpdateTime: Time of the most recent update to this course.
	// Specifying this field in a course update mask results in an
	// error.
	//
	// Read-only.
	UpdateTime string `json:"updateTime,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "AlternateLink") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AlternateLink") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Course) MarshalJSON() ([]byte, error) {
	type noMethod Course
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CourseAlias: Alternative identifier for a course.
//
// An alias uniquely identifies a course. It must be unique within one
// of the
// following scopes:
//
// * domain: A domain-scoped alias is visible to all users within the
// alias
// creator's domain and can be created only by a domain admin. A
// domain-scoped
// alias is often used when a course has an identifier external to
// Classroom.
//
// * project: A project-scoped alias is visible to any request from
// an
// application using the Developer Console project ID that created the
// alias
// and can be created by any project. A project-scoped alias is often
// used when
// an application has alternative identifiers. A random value can also
// be used
// to avoid duplicate courses in the event of transmission failures, as
// retrying
// a request will return `ALREADY_EXISTS` if a previous one has
// succeeded.
type CourseAlias struct {
	// Alias: Alias string. The format of the string indicates the desired
	// alias scoping.
	//
	// * `d:<name>` indicates a domain-scoped alias.
	//   Example: `d:math_101`
	// * `p:<name>` indicates a project-scoped alias.
	//   Example: `p:abc123`
	//
	// This field has a maximum length of 256 characters.
	Alias string `json:"alias,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Alias") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Alias") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *CourseAlias) MarshalJSON() ([]byte, error) {
	type noMethod CourseAlias
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CourseMaterial: A material attached to a course as part of a material
// set.
type CourseMaterial struct {
	// DriveFile: Google Drive file attachment.
	DriveFile *DriveFile `json:"driveFile,omitempty"`

	// Form: Google Forms attachment.
	Form *Form `json:"form,omitempty"`

	// Link: Link atatchment.
	Link *Link `json:"link,omitempty"`

	// YouTubeVideo: Youtube video attachment.
	YouTubeVideo *YouTubeVideo `json:"youTubeVideo,omitempty"`

	// ForceSendFields is a list of field names (e.g. "DriveFile") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "DriveFile") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *CourseMaterial) MarshalJSON() ([]byte, error) {
	type noMethod CourseMaterial
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CourseMaterialSet: A set of materials that appears on the "About"
// page of the course.
// These materials might include a syllabus, schedule, or other
// background
// information relating to the course as a whole.
type CourseMaterialSet struct {
	// Materials: Materials attached to this set.
	Materials []*CourseMaterial `json:"materials,omitempty"`

	// Title: Title for this set.
	Title string `json:"title,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Materials") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Materials") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *CourseMaterialSet) MarshalJSON() ([]byte, error) {
	type noMethod CourseMaterialSet
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CourseRosterChangesInfo: Information about a `Feed` with a
// `feed_type` of `COURSE_ROSTER_CHANGES`.
type CourseRosterChangesInfo struct {
	// CourseId: The `course_id` of the course to subscribe to roster
	// changes for.
	CourseId string `json:"courseId,omitempty"`

	// ForceSendFields is a list of field names (e.g. "CourseId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CourseId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *CourseRosterChangesInfo) MarshalJSON() ([]byte, error) {
	type noMethod CourseRosterChangesInfo
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CourseWork: Course work created by a teacher for students of the
// course.
type CourseWork struct {
	// AlternateLink: Absolute link to this course work in the Classroom web
	// UI.
	// This is only populated if `state` is `PUBLISHED`.
	//
	// Read-only.
	AlternateLink string `json:"alternateLink,omitempty"`

	// AssigneeMode: Assignee mode of the coursework.
	// If unspecified, the default value is `ALL_STUDENTS`.
	//
	// Possible values:
	//   "ASSIGNEE_MODE_UNSPECIFIED" - No mode specified. This is never
	// returned.
	//   "ALL_STUDENTS" - All students can see the item.
	// This is the default state.
	//   "INDIVIDUAL_STUDENTS" - A subset of the students can see the item.
	AssigneeMode string `json:"assigneeMode,omitempty"`

	// Assignment: Assignment details.
	// This is populated only when `work_type` is `ASSIGNMENT`.
	//
	// Read-only.
	Assignment *Assignment `json:"assignment,omitempty"`

	// AssociatedWithDeveloper: Whether this course work item is associated
	// with the Developer Console
	// project making the request.
	//
	// See google.classroom.Work.CreateCourseWork for
	// more
	// details.
	//
	// Read-only.
	AssociatedWithDeveloper bool `json:"associatedWithDeveloper,omitempty"`

	// CourseId: Identifier of the course.
	//
	// Read-only.
	CourseId string `json:"courseId,omitempty"`

	// CreationTime: Timestamp when this course work was
	// created.
	//
	// Read-only.
	CreationTime string `json:"creationTime,omitempty"`

	// CreatorUserId: Identifier for the user that created the
	// coursework.
	//
	// Read-only.
	CreatorUserId string `json:"creatorUserId,omitempty"`

	// Description: Optional description of this course work.
	// If set, the description must be a valid UTF-8 string containing no
	// more
	// than 30,000 characters.
	Description string `json:"description,omitempty"`

	// DueDate: Optional date, in UTC, that submissions for this this course
	// work are due.
	// This must be specified if `due_time` is specified.
	DueDate *Date `json:"dueDate,omitempty"`

	// DueTime: Optional time of day, in UTC, that submissions for this this
	// course work
	// are due.
	// This must be specified if `due_date` is specified.
	DueTime *TimeOfDay `json:"dueTime,omitempty"`

	// Id: Classroom-assigned identifier of this course work, unique per
	// course.
	//
	// Read-only.
	Id string `json:"id,omitempty"`

	// IndividualStudentsOptions: Identifiers of students with access to the
	// coursework.
	// This field is set only if `assigneeMode` is `INDIVIDUAL_STUDENTS`.
	// If the `assigneeMode` is `INDIVIDUAL_STUDENTS`, then only
	// students
	// specified in this field will be assigned the coursework.
	IndividualStudentsOptions *IndividualStudentsOptions `json:"individualStudentsOptions,omitempty"`

	// Materials: Additional materials.
	//
	// CourseWork must have no more than 20 material items.
	Materials []*Material `json:"materials,omitempty"`

	// MaxPoints: Maximum grade for this course work.
	// If zero or unspecified, this assignment is considered ungraded.
	// This must be a non-negative integer value.
	MaxPoints float64 `json:"maxPoints,omitempty"`

	// MultipleChoiceQuestion: Multiple choice question details.
	// For read operations, this field is populated only when `work_type`
	// is
	// `MULTIPLE_CHOICE_QUESTION`.
	// For write operations, this field must be specified when creating
	// course
	// work with a `work_type` of `MULTIPLE_CHOICE_QUESTION`, and it must
	// not be
	// set otherwise.
	MultipleChoiceQuestion *MultipleChoiceQuestion `json:"multipleChoiceQuestion,omitempty"`

	// ScheduledTime: Optional timestamp when this course work is scheduled
	// to be published.
	ScheduledTime string `json:"scheduledTime,omitempty"`

	// State: Status of this course work.
	// If unspecified, the default state is `DRAFT`.
	//
	// Possible values:
	//   "COURSE_WORK_STATE_UNSPECIFIED" - No state specified. This is never
	// returned.
	//   "PUBLISHED" - Status for work that has been published.
	// This is the default state.
	//   "DRAFT" - Status for work that is not yet published.
	// Work in this state is visible only to course teachers and
	// domain
	// administrators.
	//   "DELETED" - Status for work that was published but is now
	// deleted.
	// Work in this state is visible only to course teachers and
	// domain
	// administrators.
	// Work in this state is deleted after some time.
	State string `json:"state,omitempty"`

	// SubmissionModificationMode: Setting to determine when students are
	// allowed to modify submissions.
	// If unspecified, the default value is `MODIFIABLE_UNTIL_TURNED_IN`.
	//
	// Possible values:
	//   "SUBMISSION_MODIFICATION_MODE_UNSPECIFIED" - No modification mode
	// specified. This is never returned.
	//   "MODIFIABLE_UNTIL_TURNED_IN" - Submisisons can be modified before
	// being turned in.
	//   "MODIFIABLE" - Submisisons can be modified at any time.
	SubmissionModificationMode string `json:"submissionModificationMode,omitempty"`

	// Title: Title of this course work.
	// The title must be a valid UTF-8 string containing between 1 and
	// 3000
	// characters.
	Title string `json:"title,omitempty"`

	// UpdateTime: Timestamp of the most recent change to this course
	// work.
	//
	// Read-only.
	UpdateTime string `json:"updateTime,omitempty"`

	// WorkType: Type of this course work.
	//
	// The type is set when the course work is created and cannot be
	// changed.
	//
	// Possible values:
	//   "COURSE_WORK_TYPE_UNSPECIFIED" - No work type specified. This is
	// never returned.
	//   "ASSIGNMENT" - An assignment.
	//   "SHORT_ANSWER_QUESTION" - A short answer question.
	//   "MULTIPLE_CHOICE_QUESTION" - A multiple-choice question.
	WorkType string `json:"workType,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "AlternateLink") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AlternateLink") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *CourseWork) MarshalJSON() ([]byte, error) {
	type noMethod CourseWork
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

func (s *CourseWork) UnmarshalJSON(data []byte) error {
	type noMethod CourseWork
	var s1 struct {
		MaxPoints gensupport.JSONFloat64 `json:"maxPoints"`
		*noMethod
	}
	s1.noMethod = (*noMethod)(s)
	if err := json.Unmarshal(data, &s1); err != nil {
		return err
	}
	s.MaxPoints = float64(s1.MaxPoints)
	return nil
}

// Date: Represents a whole calendar date, e.g. date of birth. The time
// of day and
// time zone are either specified elsewhere or are not significant. The
// date
// is relative to the Proleptic Gregorian Calendar. The day may be 0
// to
// represent a year and month where the day is not significant, e.g.
// credit card
// expiration date. The year may be 0 to represent a month and day
// independent
// of year, e.g. anniversary date. Related types are
// google.type.TimeOfDay
// and `google.protobuf.Timestamp`.
type Date struct {
	// Day: Day of month. Must be from 1 to 31 and valid for the year and
	// month, or 0
	// if specifying a year/month where the day is not significant.
	Day int64 `json:"day,omitempty"`

	// Month: Month of year. Must be from 1 to 12.
	Month int64 `json:"month,omitempty"`

	// Year: Year of date. Must be from 1 to 9999, or 0 if specifying a date
	// without
	// a year.
	Year int64 `json:"year,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Day") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Day") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Date) MarshalJSON() ([]byte, error) {
	type noMethod Date
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// DriveFile: Representation of a Google Drive file.
type DriveFile struct {
	// AlternateLink: URL that can be used to access the Drive
	// item.
	//
	// Read-only.
	AlternateLink string `json:"alternateLink,omitempty"`

	// Id: Drive API resource ID.
	Id string `json:"id,omitempty"`

	// ThumbnailUrl: URL of a thumbnail image of the Drive item.
	//
	// Read-only.
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`

	// Title: Title of the Drive item.
	//
	// Read-only.
	Title string `json:"title,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AlternateLink") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AlternateLink") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *DriveFile) MarshalJSON() ([]byte, error) {
	type noMethod DriveFile
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// DriveFolder: Representation of a Google Drive folder.
type DriveFolder struct {
	// AlternateLink: URL that can be used to access the Drive
	// folder.
	//
	// Read-only.
	AlternateLink string `json:"alternateLink,omitempty"`

	// Id: Drive API resource ID.
	Id string `json:"id,omitempty"`

	// Title: Title of the Drive folder.
	//
	// Read-only.
	Title string `json:"title,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AlternateLink") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AlternateLink") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *DriveFolder) MarshalJSON() ([]byte, error) {
	type noMethod DriveFolder
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Empty: A generic empty message that you can re-use to avoid defining
// duplicated
// empty messages in your APIs. A typical example is to use it as the
// request
// or the response type of an API method. For instance:
//
//     service Foo {
//       rpc Bar(google.protobuf.Empty) returns
// (google.protobuf.Empty);
//     }
//
// The JSON representation for `Empty` is empty JSON object `{}`.
type Empty struct {
	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`
}

// Feed: A class of notifications that an application can register to
// receive.
// For example: "all roster changes for a domain".
type Feed struct {
	// CourseRosterChangesInfo: Information about a `Feed` with a
	// `feed_type` of `COURSE_ROSTER_CHANGES`.
	// This field must be specified if `feed_type` is
	// `COURSE_ROSTER_CHANGES`.
	CourseRosterChangesInfo *CourseRosterChangesInfo `json:"courseRosterChangesInfo,omitempty"`

	// FeedType: The type of feed.
	//
	// Possible values:
	//   "FEED_TYPE_UNSPECIFIED" - Should never be returned or provided.
	//   "DOMAIN_ROSTER_CHANGES" - All roster changes for a particular
	// domain.
	//
	// Notifications will be generated whenever a user joins or leaves a
	// course.
	//
	// No notifications will be generated when an invitation is created
	// or
	// deleted, but notifications will be generated when a user joins a
	// course
	// by accepting an invitation.
	//   "COURSE_ROSTER_CHANGES" - All roster changes for a particular
	// course.
	//
	// Notifications will be generated whenever a user joins or leaves a
	// course.
	//
	// No notifications will be generated when an invitation is created
	// or
	// deleted, but notifications will be generated when a user joins a
	// course
	// by accepting an invitation.
	FeedType string `json:"feedType,omitempty"`

	// ForceSendFields is a list of field names (e.g.
	// "CourseRosterChangesInfo") to unconditionally include in API
	// requests. By default, fields with empty values are omitted from API
	// requests. However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CourseRosterChangesInfo")
	// to include in API requests with the JSON null value. By default,
	// fields with empty values are omitted from API requests. However, any
	// field with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *Feed) MarshalJSON() ([]byte, error) {
	type noMethod Feed
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Form: Google Forms item.
type Form struct {
	// FormUrl: URL of the form.
	FormUrl string `json:"formUrl,omitempty"`

	// ResponseUrl: URL of the form responses document.
	// Only set if respsonses have been recorded and only when
	// the
	// requesting user is an editor of the form.
	//
	// Read-only.
	ResponseUrl string `json:"responseUrl,omitempty"`

	// ThumbnailUrl: URL of a thumbnail image of the Form.
	//
	// Read-only.
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`

	// Title: Title of the Form.
	//
	// Read-only.
	Title string `json:"title,omitempty"`

	// ForceSendFields is a list of field names (e.g. "FormUrl") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "FormUrl") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Form) MarshalJSON() ([]byte, error) {
	type noMethod Form
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// GlobalPermission: Global user permission description.
type GlobalPermission struct {
	// Permission: Permission value.
	//
	// Possible values:
	//   "PERMISSION_UNSPECIFIED" - No permission is specified. This is not
	// returned and is not a
	// valid value.
	//   "CREATE_COURSE" - User is permitted to create a course.
	Permission string `json:"permission,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Permission") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Permission") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *GlobalPermission) MarshalJSON() ([]byte, error) {
	type noMethod GlobalPermission
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// GradeHistory: The history of each grade on this submission.
type GradeHistory struct {
	// ActorUserId: The teacher who made the grade change.
	ActorUserId string `json:"actorUserId,omitempty"`

	// GradeChangeType: The type of grade change at this time in the
	// submission grade history.
	//
	// Possible values:
	//   "UNKNOWN_GRADE_CHANGE_TYPE" - No grade change type specified. This
	// should never be returned.
	//   "DRAFT_GRADE_POINTS_EARNED_CHANGE" - A change in the numerator of
	// the draft grade.
	//   "ASSIGNED_GRADE_POINTS_EARNED_CHANGE" - A change in the numerator
	// of the assigned grade.
	//   "MAX_POINTS_CHANGE" - A change in the denominator of the grade.
	GradeChangeType string `json:"gradeChangeType,omitempty"`

	// GradeTimestamp: When the grade of the submission was changed.
	GradeTimestamp string `json:"gradeTimestamp,omitempty"`

	// MaxPoints: The denominator of the grade at this time in the
	// submission grade
	// history.
	MaxPoints float64 `json:"maxPoints,omitempty"`

	// PointsEarned: The numerator of the grade at this time in the
	// submission grade history.
	PointsEarned float64 `json:"pointsEarned,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ActorUserId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ActorUserId") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *GradeHistory) MarshalJSON() ([]byte, error) {
	type noMethod GradeHistory
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

func (s *GradeHistory) UnmarshalJSON(data []byte) error {
	type noMethod GradeHistory
	var s1 struct {
		MaxPoints    gensupport.JSONFloat64 `json:"maxPoints"`
		PointsEarned gensupport.JSONFloat64 `json:"pointsEarned"`
		*noMethod
	}
	s1.noMethod = (*noMethod)(s)
	if err := json.Unmarshal(data, &s1); err != nil {
		return err
	}
	s.MaxPoints = float64(s1.MaxPoints)
	s.PointsEarned = float64(s1.PointsEarned)
	return nil
}

// Guardian: Association between a student and a guardian of that
// student. The guardian
// may receive information about the student's course work.
type Guardian struct {
	// GuardianId: Identifier for the guardian.
	GuardianId string `json:"guardianId,omitempty"`

	// GuardianProfile: User profile for the guardian.
	GuardianProfile *UserProfile `json:"guardianProfile,omitempty"`

	// InvitedEmailAddress: The email address to which the initial guardian
	// invitation was sent.
	// This field is only visible to domain administrators.
	InvitedEmailAddress string `json:"invitedEmailAddress,omitempty"`

	// StudentId: Identifier for the student to whom the guardian
	// relationship applies.
	StudentId string `json:"studentId,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "GuardianId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "GuardianId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Guardian) MarshalJSON() ([]byte, error) {
	type noMethod Guardian
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// GuardianInvitation: An invitation to become the guardian of a
// specified user, sent to a specified
// email address.
type GuardianInvitation struct {
	// CreationTime: The time that this invitation was created.
	//
	// Read-only.
	CreationTime string `json:"creationTime,omitempty"`

	// InvitationId: Unique identifier for this invitation.
	//
	// Read-only.
	InvitationId string `json:"invitationId,omitempty"`

	// InvitedEmailAddress: Email address that the invitation was sent
	// to.
	// This field is only visible to domain administrators.
	InvitedEmailAddress string `json:"invitedEmailAddress,omitempty"`

	// State: The state that this invitation is in.
	//
	// Possible values:
	//   "GUARDIAN_INVITATION_STATE_UNSPECIFIED" - Should never be returned.
	//   "PENDING" - The invitation is active and awaiting a response.
	//   "COMPLETE" - The invitation is no longer active. It may have been
	// accepted, declined,
	// withdrawn or it may have expired.
	State string `json:"state,omitempty"`

	// StudentId: ID of the student (in standard format)
	StudentId string `json:"studentId,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "CreationTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CreationTime") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *GuardianInvitation) MarshalJSON() ([]byte, error) {
	type noMethod GuardianInvitation
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// IndividualStudentsOptions: Assignee details about a
// coursework/announcement.
// This field is set if and only if `assigneeMode` is
// `INDIVIDUAL_STUDENTS`.
type IndividualStudentsOptions struct {
	// StudentIds: Identifiers for the students that have access to
	// the
	// coursework/announcement.
	StudentIds []string `json:"studentIds,omitempty"`

	// ForceSendFields is a list of field names (e.g. "StudentIds") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "StudentIds") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *IndividualStudentsOptions) MarshalJSON() ([]byte, error) {
	type noMethod IndividualStudentsOptions
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Invitation: An invitation to join a course.
type Invitation struct {
	// CourseId: Identifier of the course to invite the user to.
	CourseId string `json:"courseId,omitempty"`

	// Id: Identifier assigned by Classroom.
	//
	// Read-only.
	Id string `json:"id,omitempty"`

	// Role: Role to invite the user to have.
	// Must not be `COURSE_ROLE_UNSPECIFIED`.
	//
	// Possible values:
	//   "COURSE_ROLE_UNSPECIFIED" - No course role.
	//   "STUDENT" - Student in the course.
	//   "TEACHER" - Teacher of the course.
	//   "OWNER" - Owner of the course.
	Role string `json:"role,omitempty"`

	// UserId: Identifier of the invited user.
	//
	// When specified as a parameter of a request, this identifier can be
	// set to
	// one of the following:
	//
	// * the numeric identifier for the user
	// * the email address of the user
	// * the string literal "me", indicating the requesting user
	UserId string `json:"userId,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "CourseId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CourseId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Invitation) MarshalJSON() ([]byte, error) {
	type noMethod Invitation
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Link: URL item.
type Link struct {
	// ThumbnailUrl: URL of a thumbnail image of the target URL.
	//
	// Read-only.
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`

	// Title: Title of the target of the URL.
	//
	// Read-only.
	Title string `json:"title,omitempty"`

	// Url: URL to link to.
	// This must be a valid UTF-8 string containing between 1 and 2024
	// characters.
	Url string `json:"url,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ThumbnailUrl") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ThumbnailUrl") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Link) MarshalJSON() ([]byte, error) {
	type noMethod Link
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListAnnouncementsResponse: Response when listing course work.
type ListAnnouncementsResponse struct {
	// Announcements: Announcement items that match the request.
	Announcements []*Announcement `json:"announcements,omitempty"`

	// NextPageToken: Token identifying the next page of results to return.
	// If empty, no further
	// results are available.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Announcements") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Announcements") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListAnnouncementsResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListAnnouncementsResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListCourseAliasesResponse: Response when listing course aliases.
type ListCourseAliasesResponse struct {
	// Aliases: The course aliases.
	Aliases []*CourseAlias `json:"aliases,omitempty"`

	// NextPageToken: Token identifying the next page of results to return.
	// If empty, no further
	// results are available.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Aliases") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Aliases") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListCourseAliasesResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListCourseAliasesResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListCourseWorkResponse: Response when listing course work.
type ListCourseWorkResponse struct {
	// CourseWork: Course work items that match the request.
	CourseWork []*CourseWork `json:"courseWork,omitempty"`

	// NextPageToken: Token identifying the next page of results to return.
	// If empty, no further
	// results are available.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "CourseWork") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CourseWork") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListCourseWorkResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListCourseWorkResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListCoursesResponse: Response when listing courses.
type ListCoursesResponse struct {
	// Courses: Courses that match the list request.
	Courses []*Course `json:"courses,omitempty"`

	// NextPageToken: Token identifying the next page of results to return.
	// If empty, no further
	// results are available.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Courses") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Courses") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListCoursesResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListCoursesResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListGuardianInvitationsResponse: Response when listing guardian
// invitations.
type ListGuardianInvitationsResponse struct {
	// GuardianInvitations: Guardian invitations that matched the list
	// request.
	GuardianInvitations []*GuardianInvitation `json:"guardianInvitations,omitempty"`

	// NextPageToken: Token identifying the next page of results to return.
	// If empty, no further
	// results are available.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "GuardianInvitations")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "GuardianInvitations") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *ListGuardianInvitationsResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListGuardianInvitationsResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListGuardiansResponse: Response when listing guardians.
type ListGuardiansResponse struct {
	// Guardians: Guardians on this page of results that met the criteria
	// specified in
	// the request.
	Guardians []*Guardian `json:"guardians,omitempty"`

	// NextPageToken: Token identifying the next page of results to return.
	// If empty, no further
	// results are available.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Guardians") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Guardians") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListGuardiansResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListGuardiansResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListInvitationsResponse: Response when listing invitations.
type ListInvitationsResponse struct {
	// Invitations: Invitations that match the list request.
	Invitations []*Invitation `json:"invitations,omitempty"`

	// NextPageToken: Token identifying the next page of results to return.
	// If empty, no further
	// results are available.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Invitations") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Invitations") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListInvitationsResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListInvitationsResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListStudentSubmissionsResponse: Response when listing student
// submissions.
type ListStudentSubmissionsResponse struct {
	// NextPageToken: Token identifying the next page of results to return.
	// If empty, no further
	// results are available.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// StudentSubmissions: Student work that matches the request.
	StudentSubmissions []*StudentSubmission `json:"studentSubmissions,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "NextPageToken") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "NextPageToken") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListStudentSubmissionsResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListStudentSubmissionsResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListStudentsResponse: Response when listing students.
type ListStudentsResponse struct {
	// NextPageToken: Token identifying the next page of results to return.
	// If empty, no further
	// results are available.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// Students: Students who match the list request.
	Students []*Student `json:"students,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "NextPageToken") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "NextPageToken") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListStudentsResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListStudentsResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListTeachersResponse: Response when listing teachers.
type ListTeachersResponse struct {
	// NextPageToken: Token identifying the next page of results to return.
	// If empty, no further
	// results are available.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// Teachers: Teachers who match the list request.
	Teachers []*Teacher `json:"teachers,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "NextPageToken") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "NextPageToken") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListTeachersResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListTeachersResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Material: Material attached to course work.
//
// When creating attachments, setting the `form` field is not supported.
type Material struct {
	// DriveFile: Google Drive file material.
	DriveFile *SharedDriveFile `json:"driveFile,omitempty"`

	// Form: Google Forms material.
	Form *Form `json:"form,omitempty"`

	// Link: Link material. On creation, will be upgraded to a more
	// appropriate type
	// if possible, and this will be reflected in the response.
	Link *Link `json:"link,omitempty"`

	// YoutubeVideo: YouTube video material.
	YoutubeVideo *YouTubeVideo `json:"youtubeVideo,omitempty"`

	// ForceSendFields is a list of field names (e.g. "DriveFile") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "DriveFile") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Material) MarshalJSON() ([]byte, error) {
	type noMethod Material
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ModifyAnnouncementAssigneesRequest: Request to modify assignee mode
// and options of an announcement.
type ModifyAnnouncementAssigneesRequest struct {
	// AssigneeMode: Mode of the announcement describing whether it will be
	// accessible by all
	// students or specified individual students.
	//
	// Possible values:
	//   "ASSIGNEE_MODE_UNSPECIFIED" - No mode specified. This is never
	// returned.
	//   "ALL_STUDENTS" - All students can see the item.
	// This is the default state.
	//   "INDIVIDUAL_STUDENTS" - A subset of the students can see the item.
	AssigneeMode string `json:"assigneeMode,omitempty"`

	// ModifyIndividualStudentsOptions: Set which students can view or
	// cannot view the announcement.
	// Must be specified only when `assigneeMode` is `INDIVIDUAL_STUDENTS`.
	ModifyIndividualStudentsOptions *ModifyIndividualStudentsOptions `json:"modifyIndividualStudentsOptions,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AssigneeMode") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AssigneeMode") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ModifyAnnouncementAssigneesRequest) MarshalJSON() ([]byte, error) {
	type noMethod ModifyAnnouncementAssigneesRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ModifyAttachmentsRequest: Request to modify the attachments of a
// student submission.
type ModifyAttachmentsRequest struct {
	// AddAttachments: Attachments to add.
	// A student submission may not have more than 20 attachments.
	//
	// Form attachments are not supported.
	AddAttachments []*Attachment `json:"addAttachments,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AddAttachments") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AddAttachments") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *ModifyAttachmentsRequest) MarshalJSON() ([]byte, error) {
	type noMethod ModifyAttachmentsRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ModifyCourseWorkAssigneesRequest: Request to modify assignee mode and
// options of a coursework.
type ModifyCourseWorkAssigneesRequest struct {
	// AssigneeMode: Mode of the coursework describing whether it will be
	// assigned to all
	// students or specified individual students.
	//
	// Possible values:
	//   "ASSIGNEE_MODE_UNSPECIFIED" - No mode specified. This is never
	// returned.
	//   "ALL_STUDENTS" - All students can see the item.
	// This is the default state.
	//   "INDIVIDUAL_STUDENTS" - A subset of the students can see the item.
	AssigneeMode string `json:"assigneeMode,omitempty"`

	// ModifyIndividualStudentsOptions: Set which students are assigned or
	// not assigned to the coursework.
	// Must be specified only when `assigneeMode` is `INDIVIDUAL_STUDENTS`.
	ModifyIndividualStudentsOptions *ModifyIndividualStudentsOptions `json:"modifyIndividualStudentsOptions,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AssigneeMode") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AssigneeMode") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ModifyCourseWorkAssigneesRequest) MarshalJSON() ([]byte, error) {
	type noMethod ModifyCourseWorkAssigneesRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ModifyIndividualStudentsOptions: Contains fields to add or remove
// students from a course work or announcement
// where the `assigneeMode` is set to `INDIVIDUAL_STUDENTS`.
type ModifyIndividualStudentsOptions struct {
	// AddStudentIds: Ids of students to be added as having access to
	// this
	// coursework/announcement.
	AddStudentIds []string `json:"addStudentIds,omitempty"`

	// RemoveStudentIds: Ids of students to be removed from having access to
	// this
	// coursework/announcement.
	RemoveStudentIds []string `json:"removeStudentIds,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AddStudentIds") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AddStudentIds") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ModifyIndividualStudentsOptions) MarshalJSON() ([]byte, error) {
	type noMethod ModifyIndividualStudentsOptions
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// MultipleChoiceQuestion: Additional details for multiple-choice
// questions.
type MultipleChoiceQuestion struct {
	// Choices: Possible choices.
	Choices []string `json:"choices,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Choices") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Choices") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *MultipleChoiceQuestion) MarshalJSON() ([]byte, error) {
	type noMethod MultipleChoiceQuestion
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// MultipleChoiceSubmission: Student work for a multiple-choice
// question.
type MultipleChoiceSubmission struct {
	// Answer: Student's select choice.
	Answer string `json:"answer,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Answer") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Answer") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *MultipleChoiceSubmission) MarshalJSON() ([]byte, error) {
	type noMethod MultipleChoiceSubmission
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Name: Details of the user's name.
type Name struct {
	// FamilyName: The user's last name.
	//
	// Read-only.
	FamilyName string `json:"familyName,omitempty"`

	// FullName: The user's full name formed by concatenating the first and
	// last name
	// values.
	//
	// Read-only.
	FullName string `json:"fullName,omitempty"`

	// GivenName: The user's first name.
	//
	// Read-only.
	GivenName string `json:"givenName,omitempty"`

	// ForceSendFields is a list of field names (e.g. "FamilyName") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "FamilyName") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Name) MarshalJSON() ([]byte, error) {
	type noMethod Name
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ReclaimStudentSubmissionRequest: Request to reclaim a student
// submission.
type ReclaimStudentSubmissionRequest struct {
}

// Registration: An instruction to Classroom to send notifications from
// the `feed` to the
// provided `destination`.
type Registration struct {
	// CloudPubsubTopic: The Cloud Pub/Sub topic that notifications are to
	// be sent to.
	CloudPubsubTopic *CloudPubsubTopic `json:"cloudPubsubTopic,omitempty"`

	// ExpiryTime: The time until which the `Registration` is
	// effective.
	//
	// This is a read-only field assigned by the server.
	ExpiryTime string `json:"expiryTime,omitempty"`

	// Feed: Specification for the class of notifications that Classroom
	// should deliver
	// to the `destination`.
	Feed *Feed `json:"feed,omitempty"`

	// RegistrationId: A server-generated unique identifier for this
	// `Registration`.
	//
	// Read-only.
	RegistrationId string `json:"registrationId,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "CloudPubsubTopic") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CloudPubsubTopic") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *Registration) MarshalJSON() ([]byte, error) {
	type noMethod Registration
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ReturnStudentSubmissionRequest: Request to return a student
// submission.
type ReturnStudentSubmissionRequest struct {
}

// SharedDriveFile: Drive file that is used as material for course work.
type SharedDriveFile struct {
	// DriveFile: Drive file details.
	DriveFile *DriveFile `json:"driveFile,omitempty"`

	// ShareMode: Mechanism by which students access the Drive item.
	//
	// Possible values:
	//   "UNKNOWN_SHARE_MODE" - No sharing mode specified. This should never
	// be returned.
	//   "VIEW" - Students can view the shared file.
	//   "EDIT" - Students can edit the shared file.
	//   "STUDENT_COPY" - Students have a personal copy of the shared file.
	ShareMode string `json:"shareMode,omitempty"`

	// ForceSendFields is a list of field names (e.g. "DriveFile") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "DriveFile") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *SharedDriveFile) MarshalJSON() ([]byte, error) {
	type noMethod SharedDriveFile
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ShortAnswerSubmission: Student work for a short answer question.
type ShortAnswerSubmission struct {
	// Answer: Student response to a short-answer question.
	Answer string `json:"answer,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Answer") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Answer") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ShortAnswerSubmission) MarshalJSON() ([]byte, error) {
	type noMethod ShortAnswerSubmission
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// StateHistory: The history of each state this submission has been in.
type StateHistory struct {
	// ActorUserId: The teacher or student who made the change
	ActorUserId string `json:"actorUserId,omitempty"`

	// State: The workflow pipeline stage.
	//
	// Possible values:
	//   "STATE_UNSPECIFIED" - No state specified. This should never be
	// returned.
	//   "CREATED" - The Submission has been created.
	//   "TURNED_IN" - The student has turned in an assigned document, which
	// may or may not be
	// a template.
	//   "RETURNED" - The teacher has returned the assigned document to the
	// student.
	//   "RECLAIMED_BY_STUDENT" - The student turned in the assigned
	// document, and then chose to
	// "unsubmit" the assignment, giving the student control again as
	// the
	// owner.
	//   "STUDENT_EDITED_AFTER_TURN_IN" - The student edited their
	// submission after turning it in. Currently,
	// only used by Questions, when the student edits their answer.
	State string `json:"state,omitempty"`

	// StateTimestamp: When the submission entered this state.
	StateTimestamp string `json:"stateTimestamp,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ActorUserId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ActorUserId") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *StateHistory) MarshalJSON() ([]byte, error) {
	type noMethod StateHistory
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Student: Student in a course.
type Student struct {
	// CourseId: Identifier of the course.
	//
	// Read-only.
	CourseId string `json:"courseId,omitempty"`

	// Profile: Global user information for the student.
	//
	// Read-only.
	Profile *UserProfile `json:"profile,omitempty"`

	// StudentWorkFolder: Information about a Drive Folder for this
	// student's work in this course.
	// Only visible to the student and domain administrators.
	//
	// Read-only.
	StudentWorkFolder *DriveFolder `json:"studentWorkFolder,omitempty"`

	// UserId: Identifier of the user.
	//
	// When specified as a parameter of a request, this identifier can be
	// one of
	// the following:
	//
	// * the numeric identifier for the user
	// * the email address of the user
	// * the string literal "me", indicating the requesting user
	UserId string `json:"userId,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "CourseId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CourseId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Student) MarshalJSON() ([]byte, error) {
	type noMethod Student
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// StudentSubmission: Student submission for course
// work.
//
// StudentSubmission items are generated when a CourseWork item is
// created.
//
// StudentSubmissions that have never been accessed (i.e. with `state` =
// NEW)
// may not have a creation time or update time.
type StudentSubmission struct {
	// AlternateLink: Absolute link to the submission in the Classroom web
	// UI.
	//
	// Read-only.
	AlternateLink string `json:"alternateLink,omitempty"`

	// AssignedGrade: Optional grade. If unset, no grade was set.
	// This value must be non-negative. Decimal (i.e. non-integer) values
	// are
	// allowed, but will be rounded to two decimal places.
	//
	// This may be modified only by course teachers.
	AssignedGrade float64 `json:"assignedGrade,omitempty"`

	// AssignmentSubmission: Submission content when course_work_type is
	// ASSIGNMENT.
	//
	// Students can modify this content
	// using
	// google.classroom.Work.ModifyAttachments.
	AssignmentSubmission *AssignmentSubmission `json:"assignmentSubmission,omitempty"`

	// AssociatedWithDeveloper: Whether this student submission is
	// associated with the Developer Console
	// project making the request.
	//
	// See google.classroom.Work.CreateCourseWork for
	// more
	// details.
	//
	// Read-only.
	AssociatedWithDeveloper bool `json:"associatedWithDeveloper,omitempty"`

	// CourseId: Identifier of the course.
	//
	// Read-only.
	CourseId string `json:"courseId,omitempty"`

	// CourseWorkId: Identifier for the course work this corresponds
	// to.
	//
	// Read-only.
	CourseWorkId string `json:"courseWorkId,omitempty"`

	// CourseWorkType: Type of course work this submission is
	// for.
	//
	// Read-only.
	//
	// Possible values:
	//   "COURSE_WORK_TYPE_UNSPECIFIED" - No work type specified. This is
	// never returned.
	//   "ASSIGNMENT" - An assignment.
	//   "SHORT_ANSWER_QUESTION" - A short answer question.
	//   "MULTIPLE_CHOICE_QUESTION" - A multiple-choice question.
	CourseWorkType string `json:"courseWorkType,omitempty"`

	// CreationTime: Creation time of this submission.
	// This may be unset if the student has not accessed this
	// item.
	//
	// Read-only.
	CreationTime string `json:"creationTime,omitempty"`

	// DraftGrade: Optional pending grade. If unset, no grade was set.
	// This value must be non-negative. Decimal (i.e. non-integer) values
	// are
	// allowed, but will be rounded to two decimal places.
	//
	// This is only visible to and modifiable by course teachers.
	DraftGrade float64 `json:"draftGrade,omitempty"`

	// Id: Classroom-assigned Identifier for the student submission.
	// This is unique among submissions for the relevant course
	// work.
	//
	// Read-only.
	Id string `json:"id,omitempty"`

	// Late: Whether this submission is late.
	//
	// Read-only.
	Late bool `json:"late,omitempty"`

	// MultipleChoiceSubmission: Submission content when course_work_type is
	// MULTIPLE_CHOICE_QUESTION.
	MultipleChoiceSubmission *MultipleChoiceSubmission `json:"multipleChoiceSubmission,omitempty"`

	// ShortAnswerSubmission: Submission content when course_work_type is
	// SHORT_ANSWER_QUESTION.
	ShortAnswerSubmission *ShortAnswerSubmission `json:"shortAnswerSubmission,omitempty"`

	// State: State of this submission.
	//
	// Read-only.
	//
	// Possible values:
	//   "SUBMISSION_STATE_UNSPECIFIED" - No state specified. This should
	// never be returned.
	//   "NEW" - The student has never accessed this submission. Attachments
	// are not
	// returned and timestamps is not set.
	//   "CREATED" - Has been created.
	//   "TURNED_IN" - Has been turned in to the teacher.
	//   "RETURNED" - Has been returned to the student.
	//   "RECLAIMED_BY_STUDENT" - Student chose to "unsubmit" the
	// assignment.
	State string `json:"state,omitempty"`

	// SubmissionHistory: The history of the submission (includes state and
	// grade histories).
	//
	// Read-only.
	SubmissionHistory []*SubmissionHistory `json:"submissionHistory,omitempty"`

	// UpdateTime: Last update time of this submission.
	// This may be unset if the student has not accessed this
	// item.
	//
	// Read-only.
	UpdateTime string `json:"updateTime,omitempty"`

	// UserId: Identifier for the student that owns this
	// submission.
	//
	// Read-only.
	UserId string `json:"userId,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "AlternateLink") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AlternateLink") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *StudentSubmission) MarshalJSON() ([]byte, error) {
	type noMethod StudentSubmission
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

func (s *StudentSubmission) UnmarshalJSON(data []byte) error {
	type noMethod StudentSubmission
	var s1 struct {
		AssignedGrade gensupport.JSONFloat64 `json:"assignedGrade"`
		DraftGrade    gensupport.JSONFloat64 `json:"draftGrade"`
		*noMethod
	}
	s1.noMethod = (*noMethod)(s)
	if err := json.Unmarshal(data, &s1); err != nil {
		return err
	}
	s.AssignedGrade = float64(s1.AssignedGrade)
	s.DraftGrade = float64(s1.DraftGrade)
	return nil
}

// SubmissionHistory: The history of the submission. This currently
// includes state and grade
// histories.
type SubmissionHistory struct {
	// GradeHistory: The grade history information of the submission, if
	// present.
	GradeHistory *GradeHistory `json:"gradeHistory,omitempty"`

	// StateHistory: The state history information of the submission, if
	// present.
	StateHistory *StateHistory `json:"stateHistory,omitempty"`

	// ForceSendFields is a list of field names (e.g. "GradeHistory") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "GradeHistory") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *SubmissionHistory) MarshalJSON() ([]byte, error) {
	type noMethod SubmissionHistory
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Teacher: Teacher of a course.
type Teacher struct {
	// CourseId: Identifier of the course.
	//
	// Read-only.
	CourseId string `json:"courseId,omitempty"`

	// Profile: Global user information for the teacher.
	//
	// Read-only.
	Profile *UserProfile `json:"profile,omitempty"`

	// UserId: Identifier of the user.
	//
	// When specified as a parameter of a request, this identifier can be
	// one of
	// the following:
	//
	// * the numeric identifier for the user
	// * the email address of the user
	// * the string literal "me", indicating the requesting user
	UserId string `json:"userId,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "CourseId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CourseId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Teacher) MarshalJSON() ([]byte, error) {
	type noMethod Teacher
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// TimeOfDay: Represents a time of day. The date and time zone are
// either not significant
// or are specified elsewhere. An API may choose to allow leap seconds.
// Related
// types are google.type.Date and `google.protobuf.Timestamp`.
type TimeOfDay struct {
	// Hours: Hours of day in 24 hour format. Should be from 0 to 23. An API
	// may choose
	// to allow the value "24:00:00" for scenarios like business closing
	// time.
	Hours int64 `json:"hours,omitempty"`

	// Minutes: Minutes of hour of day. Must be from 0 to 59.
	Minutes int64 `json:"minutes,omitempty"`

	// Nanos: Fractions of seconds in nanoseconds. Must be from 0 to
	// 999,999,999.
	Nanos int64 `json:"nanos,omitempty"`

	// Seconds: Seconds of minutes of the time. Must normally be from 0 to
	// 59. An API may
	// allow the value 60 if it allows leap-seconds.
	Seconds int64 `json:"seconds,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Hours") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Hours") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *TimeOfDay) MarshalJSON() ([]byte, error) {
	type noMethod TimeOfDay
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// TurnInStudentSubmissionRequest: Request to turn in a student
// submission.
type TurnInStudentSubmissionRequest struct {
}

// UserProfile: Global information for a user.
type UserProfile struct {
	// EmailAddress: Email address of the user.
	//
	// Read-only.
	EmailAddress string `json:"emailAddress,omitempty"`

	// Id: Identifier of the user.
	//
	// Read-only.
	Id string `json:"id,omitempty"`

	// Name: Name of the user.
	//
	// Read-only.
	Name *Name `json:"name,omitempty"`

	// Permissions: Global permissions of the user.
	//
	// Read-only.
	Permissions []*GlobalPermission `json:"permissions,omitempty"`

	// PhotoUrl: URL of user's profile photo.
	//
	// Read-only.
	PhotoUrl string `json:"photoUrl,omitempty"`

	// VerifiedTeacher: Represents whether a G Suite for Education user's
	// domain administrator has
	// explicitly verified them as being a teacher. If the user is not a
	// member of
	// a G Suite for Education domain, than this field will always be
	// false.
	//
	// Read-only
	VerifiedTeacher bool `json:"verifiedTeacher,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "EmailAddress") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "EmailAddress") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *UserProfile) MarshalJSON() ([]byte, error) {
	type noMethod UserProfile
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// YouTubeVideo: YouTube video item.
type YouTubeVideo struct {
	// AlternateLink: URL that can be used to view the YouTube
	// video.
	//
	// Read-only.
	AlternateLink string `json:"alternateLink,omitempty"`

	// Id: YouTube API resource ID.
	Id string `json:"id,omitempty"`

	// ThumbnailUrl: URL of a thumbnail image of the YouTube
	// video.
	//
	// Read-only.
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`

	// Title: Title of the YouTube video.
	//
	// Read-only.
	Title string `json:"title,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AlternateLink") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AlternateLink") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *YouTubeVideo) MarshalJSON() ([]byte, error) {
	type noMethod YouTubeVideo
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// method id "classroom.courses.create":

type CoursesCreateCall struct {
	s          *Service
	course     *Course
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Create: Creates a course.
//
// The user specified in `ownerId` is the owner of the created
// course
// and added as a teacher.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// create
// courses or for access errors.
// * `NOT_FOUND` if the primary teacher is not a valid user.
// * `FAILED_PRECONDITION` if the course owner's account is disabled or
// for
// the following request errors:
//     * UserGroupsMembershipLimitReached
// * `ALREADY_EXISTS` if an alias was specified in the `id` and
// already exists.
func (r *CoursesService) Create(course *Course) *CoursesCreateCall {
	c := &CoursesCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.course = course
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesCreateCall) Fields(s ...googleapi.Field) *CoursesCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesCreateCall) Context(ctx context.Context) *CoursesCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.course)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.create" call.
// Exactly one of *Course or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Course.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesCreateCall) Do(opts ...googleapi.CallOption) (*Course, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Course{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Creates a course.\n\nThe user specified in `ownerId` is the owner of the created course\nand added as a teacher.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to create\ncourses or for access errors.\n* `NOT_FOUND` if the primary teacher is not a valid user.\n* `FAILED_PRECONDITION` if the course owner's account is disabled or for\nthe following request errors:\n    * UserGroupsMembershipLimitReached\n* `ALREADY_EXISTS` if an alias was specified in the `id` and\nalready exists.",
	//   "flatPath": "v1/courses",
	//   "httpMethod": "POST",
	//   "id": "classroom.courses.create",
	//   "parameterOrder": [],
	//   "parameters": {},
	//   "path": "v1/courses",
	//   "request": {
	//     "$ref": "Course"
	//   },
	//   "response": {
	//     "$ref": "Course"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.courses"
	//   ]
	// }

}

// method id "classroom.courses.delete":

type CoursesDeleteCall struct {
	s          *Service
	id         string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes a course.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// delete the
// requested course or for access errors.
// * `NOT_FOUND` if no course exists with the requested ID.
func (r *CoursesService) Delete(id string) *CoursesDeleteCall {
	c := &CoursesDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.id = id
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesDeleteCall) Fields(s ...googleapi.Field) *CoursesDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesDeleteCall) Context(ctx context.Context) *CoursesDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{id}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"id": c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.delete" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesDeleteCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Deletes a course.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to delete the\nrequested course or for access errors.\n* `NOT_FOUND` if no course exists with the requested ID.",
	//   "flatPath": "v1/courses/{id}",
	//   "httpMethod": "DELETE",
	//   "id": "classroom.courses.delete",
	//   "parameterOrder": [
	//     "id"
	//   ],
	//   "parameters": {
	//     "id": {
	//       "description": "Identifier of the course to delete.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{id}",
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.courses"
	//   ]
	// }

}

// method id "classroom.courses.get":

type CoursesGetCall struct {
	s            *Service
	id           string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Returns a course.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access the
// requested course or for access errors.
// * `NOT_FOUND` if no course exists with the requested ID.
func (r *CoursesService) Get(id string) *CoursesGetCall {
	c := &CoursesGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.id = id
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesGetCall) Fields(s ...googleapi.Field) *CoursesGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *CoursesGetCall) IfNoneMatch(entityTag string) *CoursesGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesGetCall) Context(ctx context.Context) *CoursesGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesGetCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{id}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"id": c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.get" call.
// Exactly one of *Course or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Course.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesGetCall) Do(opts ...googleapi.CallOption) (*Course, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Course{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a course.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access the\nrequested course or for access errors.\n* `NOT_FOUND` if no course exists with the requested ID.",
	//   "flatPath": "v1/courses/{id}",
	//   "httpMethod": "GET",
	//   "id": "classroom.courses.get",
	//   "parameterOrder": [
	//     "id"
	//   ],
	//   "parameters": {
	//     "id": {
	//       "description": "Identifier of the course to return.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{id}",
	//   "response": {
	//     "$ref": "Course"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.courses",
	//     "https://www.googleapis.com/auth/classroom.courses.readonly"
	//   ]
	// }

}

// method id "classroom.courses.list":

type CoursesListCall struct {
	s            *Service
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Returns a list of courses that the requesting user is permitted
// to view,
// restricted to those that match the request. Returned courses are
// ordered by
// creation time, with the most recently created coming first.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` for access errors.
// * `INVALID_ARGUMENT` if the query argument is malformed.
// * `NOT_FOUND` if any users specified in the query arguments do not
// exist.
func (r *CoursesService) List() *CoursesListCall {
	c := &CoursesListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	return c
}

// CourseStates sets the optional parameter "courseStates": Restricts
// returned courses to those in one of the specified states
// The default value is ACTIVE, ARCHIVED, PROVISIONED, DECLINED.
//
// Possible values:
//   "COURSE_STATE_UNSPECIFIED"
//   "ACTIVE"
//   "ARCHIVED"
//   "PROVISIONED"
//   "DECLINED"
//   "SUSPENDED"
func (c *CoursesListCall) CourseStates(courseStates ...string) *CoursesListCall {
	c.urlParams_.SetMulti("courseStates", append([]string{}, courseStates...))
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum number of
// items to return. Zero or unspecified indicates that the
// server may assign a maximum.
//
// The server may return fewer than the specified number of results.
func (c *CoursesListCall) PageSize(pageSize int64) *CoursesListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken":
// nextPageToken
// value returned from a previous
// list call,
// indicating that the subsequent page of results should be
// returned.
//
// The list request must be
// otherwise identical to the one that resulted in this token.
func (c *CoursesListCall) PageToken(pageToken string) *CoursesListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// StudentId sets the optional parameter "studentId": Restricts returned
// courses to those having a student with the specified
// identifier. The identifier can be one of the following:
//
// * the numeric identifier for the user
// * the email address of the user
// * the string literal "me", indicating the requesting user
func (c *CoursesListCall) StudentId(studentId string) *CoursesListCall {
	c.urlParams_.Set("studentId", studentId)
	return c
}

// TeacherId sets the optional parameter "teacherId": Restricts returned
// courses to those having a teacher with the specified
// identifier. The identifier can be one of the following:
//
// * the numeric identifier for the user
// * the email address of the user
// * the string literal "me", indicating the requesting user
func (c *CoursesListCall) TeacherId(teacherId string) *CoursesListCall {
	c.urlParams_.Set("teacherId", teacherId)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesListCall) Fields(s ...googleapi.Field) *CoursesListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *CoursesListCall) IfNoneMatch(entityTag string) *CoursesListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesListCall) Context(ctx context.Context) *CoursesListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesListCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.list" call.
// Exactly one of *ListCoursesResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListCoursesResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *CoursesListCall) Do(opts ...googleapi.CallOption) (*ListCoursesResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListCoursesResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a list of courses that the requesting user is permitted to view,\nrestricted to those that match the request. Returned courses are ordered by\ncreation time, with the most recently created coming first.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` for access errors.\n* `INVALID_ARGUMENT` if the query argument is malformed.\n* `NOT_FOUND` if any users specified in the query arguments do not exist.",
	//   "flatPath": "v1/courses",
	//   "httpMethod": "GET",
	//   "id": "classroom.courses.list",
	//   "parameterOrder": [],
	//   "parameters": {
	//     "courseStates": {
	//       "description": "Restricts returned courses to those in one of the specified states\nThe default value is ACTIVE, ARCHIVED, PROVISIONED, DECLINED.",
	//       "enum": [
	//         "COURSE_STATE_UNSPECIFIED",
	//         "ACTIVE",
	//         "ARCHIVED",
	//         "PROVISIONED",
	//         "DECLINED",
	//         "SUSPENDED"
	//       ],
	//       "location": "query",
	//       "repeated": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum number of items to return. Zero or unspecified indicates that the\nserver may assign a maximum.\n\nThe server may return fewer than the specified number of results.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "nextPageToken\nvalue returned from a previous\nlist call,\nindicating that the subsequent page of results should be returned.\n\nThe list request must be\notherwise identical to the one that resulted in this token.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "studentId": {
	//       "description": "Restricts returned courses to those having a student with the specified\nidentifier. The identifier can be one of the following:\n\n* the numeric identifier for the user\n* the email address of the user\n* the string literal `\"me\"`, indicating the requesting user",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "teacherId": {
	//       "description": "Restricts returned courses to those having a teacher with the specified\nidentifier. The identifier can be one of the following:\n\n* the numeric identifier for the user\n* the email address of the user\n* the string literal `\"me\"`, indicating the requesting user",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses",
	//   "response": {
	//     "$ref": "ListCoursesResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.courses",
	//     "https://www.googleapis.com/auth/classroom.courses.readonly"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *CoursesListCall) Pages(ctx context.Context, f func(*ListCoursesResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "classroom.courses.patch":

type CoursesPatchCall struct {
	s          *Service
	id         string
	course     *Course
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Patch: Updates one or more fields in a course.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// modify the
// requested course or for access errors.
// * `NOT_FOUND` if no course exists with the requested ID.
// * `INVALID_ARGUMENT` if invalid fields are specified in the update
// mask or
// if no update mask is supplied.
// * `FAILED_PRECONDITION` for the following request errors:
//     * CourseNotModifiable
func (r *CoursesService) Patch(id string, course *Course) *CoursesPatchCall {
	c := &CoursesPatchCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.id = id
	c.course = course
	return c
}

// UpdateMask sets the optional parameter "updateMask": Mask that
// identifies which fields on the course to update.
// This field is required to do an update. The update will fail if
// invalid
// fields are specified. The following fields are valid:
//
// * `name`
// * `section`
// * `descriptionHeading`
// * `description`
// * `room`
// * `courseState`
// * `ownerId`
//
// Note: patches to ownerId are treated as being effective immediately,
// but in
// practice it may take some time for the ownership transfer of all
// affected
// resources to complete.
//
// When set in a query parameter, this field should be specified
// as
//
// `updateMask=<field1>,<field2>,...`
func (c *CoursesPatchCall) UpdateMask(updateMask string) *CoursesPatchCall {
	c.urlParams_.Set("updateMask", updateMask)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesPatchCall) Fields(s ...googleapi.Field) *CoursesPatchCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesPatchCall) Context(ctx context.Context) *CoursesPatchCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesPatchCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesPatchCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.course)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{id}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"id": c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.patch" call.
// Exactly one of *Course or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Course.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesPatchCall) Do(opts ...googleapi.CallOption) (*Course, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Course{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates one or more fields in a course.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to modify the\nrequested course or for access errors.\n* `NOT_FOUND` if no course exists with the requested ID.\n* `INVALID_ARGUMENT` if invalid fields are specified in the update mask or\nif no update mask is supplied.\n* `FAILED_PRECONDITION` for the following request errors:\n    * CourseNotModifiable",
	//   "flatPath": "v1/courses/{id}",
	//   "httpMethod": "PATCH",
	//   "id": "classroom.courses.patch",
	//   "parameterOrder": [
	//     "id"
	//   ],
	//   "parameters": {
	//     "id": {
	//       "description": "Identifier of the course to update.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "updateMask": {
	//       "description": "Mask that identifies which fields on the course to update.\nThis field is required to do an update. The update will fail if invalid\nfields are specified. The following fields are valid:\n\n* `name`\n* `section`\n* `descriptionHeading`\n* `description`\n* `room`\n* `courseState`\n* `ownerId`\n\nNote: patches to ownerId are treated as being effective immediately, but in\npractice it may take some time for the ownership transfer of all affected\nresources to complete.\n\nWhen set in a query parameter, this field should be specified as\n\n`updateMask=\u003cfield1\u003e,\u003cfield2\u003e,...`",
	//       "format": "google-fieldmask",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{id}",
	//   "request": {
	//     "$ref": "Course"
	//   },
	//   "response": {
	//     "$ref": "Course"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.courses"
	//   ]
	// }

}

// method id "classroom.courses.update":

type CoursesUpdateCall struct {
	s          *Service
	id         string
	course     *Course
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Update: Updates a course.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// modify the
// requested course or for access errors.
// * `NOT_FOUND` if no course exists with the requested ID.
// * `FAILED_PRECONDITION` for the following request errors:
//     * CourseNotModifiable
func (r *CoursesService) Update(id string, course *Course) *CoursesUpdateCall {
	c := &CoursesUpdateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.id = id
	c.course = course
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesUpdateCall) Fields(s ...googleapi.Field) *CoursesUpdateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesUpdateCall) Context(ctx context.Context) *CoursesUpdateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesUpdateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesUpdateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.course)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{id}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PUT", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"id": c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.update" call.
// Exactly one of *Course or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Course.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesUpdateCall) Do(opts ...googleapi.CallOption) (*Course, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Course{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates a course.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to modify the\nrequested course or for access errors.\n* `NOT_FOUND` if no course exists with the requested ID.\n* `FAILED_PRECONDITION` for the following request errors:\n    * CourseNotModifiable",
	//   "flatPath": "v1/courses/{id}",
	//   "httpMethod": "PUT",
	//   "id": "classroom.courses.update",
	//   "parameterOrder": [
	//     "id"
	//   ],
	//   "parameters": {
	//     "id": {
	//       "description": "Identifier of the course to update.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{id}",
	//   "request": {
	//     "$ref": "Course"
	//   },
	//   "response": {
	//     "$ref": "Course"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.courses"
	//   ]
	// }

}

// method id "classroom.courses.aliases.create":

type CoursesAliasesCreateCall struct {
	s           *Service
	courseId    string
	coursealias *CourseAlias
	urlParams_  gensupport.URLParams
	ctx_        context.Context
	header_     http.Header
}

// Create: Creates an alias for a course.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// create the
// alias or for access errors.
// * `NOT_FOUND` if the course does not exist.
// * `ALREADY_EXISTS` if the alias already exists.
// * `FAILED_PRECONDITION` if the alias requested does not make sense
// for the
//   requesting user or course (for example, if a user not in a domain
//   attempts to access a domain-scoped alias).
func (r *CoursesAliasesService) Create(courseId string, coursealias *CourseAlias) *CoursesAliasesCreateCall {
	c := &CoursesAliasesCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.coursealias = coursealias
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesAliasesCreateCall) Fields(s ...googleapi.Field) *CoursesAliasesCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesAliasesCreateCall) Context(ctx context.Context) *CoursesAliasesCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesAliasesCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesAliasesCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.coursealias)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/aliases")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.aliases.create" call.
// Exactly one of *CourseAlias or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *CourseAlias.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *CoursesAliasesCreateCall) Do(opts ...googleapi.CallOption) (*CourseAlias, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &CourseAlias{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Creates an alias for a course.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to create the\nalias or for access errors.\n* `NOT_FOUND` if the course does not exist.\n* `ALREADY_EXISTS` if the alias already exists.\n* `FAILED_PRECONDITION` if the alias requested does not make sense for the\n  requesting user or course (for example, if a user not in a domain\n  attempts to access a domain-scoped alias).",
	//   "flatPath": "v1/courses/{courseId}/aliases",
	//   "httpMethod": "POST",
	//   "id": "classroom.courses.aliases.create",
	//   "parameterOrder": [
	//     "courseId"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course to alias.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/aliases",
	//   "request": {
	//     "$ref": "CourseAlias"
	//   },
	//   "response": {
	//     "$ref": "CourseAlias"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.courses"
	//   ]
	// }

}

// method id "classroom.courses.aliases.delete":

type CoursesAliasesDeleteCall struct {
	s          *Service
	courseId   string
	aliasid    string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes an alias of a course.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// remove the
// alias or for access errors.
// * `NOT_FOUND` if the alias does not exist.
// * `FAILED_PRECONDITION` if the alias requested does not make sense
// for the
//   requesting user or course (for example, if a user not in a domain
//   attempts to delete a domain-scoped alias).
func (r *CoursesAliasesService) Delete(courseId string, aliasid string) *CoursesAliasesDeleteCall {
	c := &CoursesAliasesDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.aliasid = aliasid
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesAliasesDeleteCall) Fields(s ...googleapi.Field) *CoursesAliasesDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesAliasesDeleteCall) Context(ctx context.Context) *CoursesAliasesDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesAliasesDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesAliasesDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/aliases/{alias}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
		"alias":    c.aliasid,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.aliases.delete" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesAliasesDeleteCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Deletes an alias of a course.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to remove the\nalias or for access errors.\n* `NOT_FOUND` if the alias does not exist.\n* `FAILED_PRECONDITION` if the alias requested does not make sense for the\n  requesting user or course (for example, if a user not in a domain\n  attempts to delete a domain-scoped alias).",
	//   "flatPath": "v1/courses/{courseId}/aliases/{alias}",
	//   "httpMethod": "DELETE",
	//   "id": "classroom.courses.aliases.delete",
	//   "parameterOrder": [
	//     "courseId",
	//     "alias"
	//   ],
	//   "parameters": {
	//     "alias": {
	//       "description": "Alias to delete.\nThis may not be the Classroom-assigned identifier.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "courseId": {
	//       "description": "Identifier of the course whose alias should be deleted.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/aliases/{alias}",
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.courses"
	//   ]
	// }

}

// method id "classroom.courses.aliases.list":

type CoursesAliasesListCall struct {
	s            *Service
	courseId     string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Returns a list of aliases for a course.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access the
// course or for access errors.
// * `NOT_FOUND` if the course does not exist.
func (r *CoursesAliasesService) List(courseId string) *CoursesAliasesListCall {
	c := &CoursesAliasesListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum number of
// items to return. Zero or unspecified indicates that the
// server may assign a maximum.
//
// The server may return fewer than the specified number of results.
func (c *CoursesAliasesListCall) PageSize(pageSize int64) *CoursesAliasesListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken":
// nextPageToken
// value returned from a previous
// list call,
// indicating that the subsequent page of results should be
// returned.
//
// The list request
// must be otherwise identical to the one that resulted in this token.
func (c *CoursesAliasesListCall) PageToken(pageToken string) *CoursesAliasesListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesAliasesListCall) Fields(s ...googleapi.Field) *CoursesAliasesListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *CoursesAliasesListCall) IfNoneMatch(entityTag string) *CoursesAliasesListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesAliasesListCall) Context(ctx context.Context) *CoursesAliasesListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesAliasesListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesAliasesListCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/aliases")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.aliases.list" call.
// Exactly one of *ListCourseAliasesResponse or error will be non-nil.
// Any non-2xx status code is an error. Response headers are in either
// *ListCourseAliasesResponse.ServerResponse.Header or (if a response
// was returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *CoursesAliasesListCall) Do(opts ...googleapi.CallOption) (*ListCourseAliasesResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListCourseAliasesResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a list of aliases for a course.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access the\ncourse or for access errors.\n* `NOT_FOUND` if the course does not exist.",
	//   "flatPath": "v1/courses/{courseId}/aliases",
	//   "httpMethod": "GET",
	//   "id": "classroom.courses.aliases.list",
	//   "parameterOrder": [
	//     "courseId"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "The identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum number of items to return. Zero or unspecified indicates that the\nserver may assign a maximum.\n\nThe server may return fewer than the specified number of results.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "nextPageToken\nvalue returned from a previous\nlist call,\nindicating that the subsequent page of results should be returned.\n\nThe list request\nmust be otherwise identical to the one that resulted in this token.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/aliases",
	//   "response": {
	//     "$ref": "ListCourseAliasesResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.courses",
	//     "https://www.googleapis.com/auth/classroom.courses.readonly"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *CoursesAliasesListCall) Pages(ctx context.Context, f func(*ListCourseAliasesResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "classroom.courses.announcements.create":

type CoursesAnnouncementsCreateCall struct {
	s            *Service
	courseId     string
	announcement *Announcement
	urlParams_   gensupport.URLParams
	ctx_         context.Context
	header_      http.Header
}

// Create: Creates an announcement.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access the
// requested course, create announcements in the requested course, share
// a
// Drive attachment, or for access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course does not exist.
// * `FAILED_PRECONDITION` for the following request error:
//     * AttachmentNotVisible
func (r *CoursesAnnouncementsService) Create(courseId string, announcement *Announcement) *CoursesAnnouncementsCreateCall {
	c := &CoursesAnnouncementsCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.announcement = announcement
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesAnnouncementsCreateCall) Fields(s ...googleapi.Field) *CoursesAnnouncementsCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesAnnouncementsCreateCall) Context(ctx context.Context) *CoursesAnnouncementsCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesAnnouncementsCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesAnnouncementsCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.announcement)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/announcements")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.announcements.create" call.
// Exactly one of *Announcement or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Announcement.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *CoursesAnnouncementsCreateCall) Do(opts ...googleapi.CallOption) (*Announcement, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Announcement{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Creates an announcement.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access the\nrequested course, create announcements in the requested course, share a\nDrive attachment, or for access errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course does not exist.\n* `FAILED_PRECONDITION` for the following request error:\n    * AttachmentNotVisible",
	//   "flatPath": "v1/courses/{courseId}/announcements",
	//   "httpMethod": "POST",
	//   "id": "classroom.courses.announcements.create",
	//   "parameterOrder": [
	//     "courseId"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/announcements",
	//   "request": {
	//     "$ref": "Announcement"
	//   },
	//   "response": {
	//     "$ref": "Announcement"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.announcements"
	//   ]
	// }

}

// method id "classroom.courses.announcements.delete":

type CoursesAnnouncementsDeleteCall struct {
	s          *Service
	courseId   string
	id         string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes an announcement.
//
// This request must be made by the Developer Console project of
// the
// [OAuth client ID](https://support.google.com/cloud/answer/6158849)
// used to
// create the corresponding announcement item.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting developer project did not
// create
// the corresponding announcement, if the requesting user is not
// permitted
// to delete the requested course or for access errors.
// * `FAILED_PRECONDITION` if the requested announcement has already
// been
// deleted.
// * `NOT_FOUND` if no course exists with the requested ID.
func (r *CoursesAnnouncementsService) Delete(courseId string, id string) *CoursesAnnouncementsDeleteCall {
	c := &CoursesAnnouncementsDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.id = id
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesAnnouncementsDeleteCall) Fields(s ...googleapi.Field) *CoursesAnnouncementsDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesAnnouncementsDeleteCall) Context(ctx context.Context) *CoursesAnnouncementsDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesAnnouncementsDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesAnnouncementsDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/announcements/{id}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
		"id":       c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.announcements.delete" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesAnnouncementsDeleteCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Deletes an announcement.\n\nThis request must be made by the Developer Console project of the\n[OAuth client ID](https://support.google.com/cloud/answer/6158849) used to\ncreate the corresponding announcement item.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting developer project did not create\nthe corresponding announcement, if the requesting user is not permitted\nto delete the requested course or for access errors.\n* `FAILED_PRECONDITION` if the requested announcement has already been\ndeleted.\n* `NOT_FOUND` if no course exists with the requested ID.",
	//   "flatPath": "v1/courses/{courseId}/announcements/{id}",
	//   "httpMethod": "DELETE",
	//   "id": "classroom.courses.announcements.delete",
	//   "parameterOrder": [
	//     "courseId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "description": "Identifier of the announcement to delete.\nThis identifier is a Classroom-assigned identifier.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/announcements/{id}",
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.announcements"
	//   ]
	// }

}

// method id "classroom.courses.announcements.get":

type CoursesAnnouncementsGetCall struct {
	s            *Service
	courseId     string
	id           string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Returns an announcement.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access the
// requested course or announcement, or for access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course or announcement does not exist.
func (r *CoursesAnnouncementsService) Get(courseId string, id string) *CoursesAnnouncementsGetCall {
	c := &CoursesAnnouncementsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.id = id
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesAnnouncementsGetCall) Fields(s ...googleapi.Field) *CoursesAnnouncementsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *CoursesAnnouncementsGetCall) IfNoneMatch(entityTag string) *CoursesAnnouncementsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesAnnouncementsGetCall) Context(ctx context.Context) *CoursesAnnouncementsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesAnnouncementsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesAnnouncementsGetCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/announcements/{id}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
		"id":       c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.announcements.get" call.
// Exactly one of *Announcement or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Announcement.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *CoursesAnnouncementsGetCall) Do(opts ...googleapi.CallOption) (*Announcement, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Announcement{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns an announcement.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access the\nrequested course or announcement, or for access errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course or announcement does not exist.",
	//   "flatPath": "v1/courses/{courseId}/announcements/{id}",
	//   "httpMethod": "GET",
	//   "id": "classroom.courses.announcements.get",
	//   "parameterOrder": [
	//     "courseId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "description": "Identifier of the announcement.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/announcements/{id}",
	//   "response": {
	//     "$ref": "Announcement"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.announcements",
	//     "https://www.googleapis.com/auth/classroom.announcements.readonly"
	//   ]
	// }

}

// method id "classroom.courses.announcements.list":

type CoursesAnnouncementsListCall struct {
	s            *Service
	courseId     string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Returns a list of announcements that the requester is permitted
// to view.
//
// Course students may only view `PUBLISHED` announcements. Course
// teachers
// and domain administrators may view all announcements.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access
// the requested course or for access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course does not exist.
func (r *CoursesAnnouncementsService) List(courseId string) *CoursesAnnouncementsListCall {
	c := &CoursesAnnouncementsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	return c
}

// AnnouncementStates sets the optional parameter "announcementStates":
// Restriction on the `state` of announcements returned.
// If this argument is left unspecified, the default value is
// `PUBLISHED`.
//
// Possible values:
//   "ANNOUNCEMENT_STATE_UNSPECIFIED"
//   "PUBLISHED"
//   "DRAFT"
//   "DELETED"
func (c *CoursesAnnouncementsListCall) AnnouncementStates(announcementStates ...string) *CoursesAnnouncementsListCall {
	c.urlParams_.SetMulti("announcementStates", append([]string{}, announcementStates...))
	return c
}

// OrderBy sets the optional parameter "orderBy": Optional sort ordering
// for results. A comma-separated list of fields with
// an optional sort direction keyword. Supported field is
// `updateTime`.
// Supported direction keywords are `asc` and `desc`.
// If not specified, `updateTime desc` is the default
// behavior.
// Examples: `updateTime asc`, `updateTime`
func (c *CoursesAnnouncementsListCall) OrderBy(orderBy string) *CoursesAnnouncementsListCall {
	c.urlParams_.Set("orderBy", orderBy)
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum number of
// items to return. Zero or unspecified indicates that the
// server may assign a maximum.
//
// The server may return fewer than the specified number of results.
func (c *CoursesAnnouncementsListCall) PageSize(pageSize int64) *CoursesAnnouncementsListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken":
// nextPageToken
// value returned from a previous
// list call,
// indicating that the subsequent page of results should be
// returned.
//
// The list request
// must be otherwise identical to the one that resulted in this token.
func (c *CoursesAnnouncementsListCall) PageToken(pageToken string) *CoursesAnnouncementsListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesAnnouncementsListCall) Fields(s ...googleapi.Field) *CoursesAnnouncementsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *CoursesAnnouncementsListCall) IfNoneMatch(entityTag string) *CoursesAnnouncementsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesAnnouncementsListCall) Context(ctx context.Context) *CoursesAnnouncementsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesAnnouncementsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesAnnouncementsListCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/announcements")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.announcements.list" call.
// Exactly one of *ListAnnouncementsResponse or error will be non-nil.
// Any non-2xx status code is an error. Response headers are in either
// *ListAnnouncementsResponse.ServerResponse.Header or (if a response
// was returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *CoursesAnnouncementsListCall) Do(opts ...googleapi.CallOption) (*ListAnnouncementsResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListAnnouncementsResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a list of announcements that the requester is permitted to view.\n\nCourse students may only view `PUBLISHED` announcements. Course teachers\nand domain administrators may view all announcements.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access\nthe requested course or for access errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course does not exist.",
	//   "flatPath": "v1/courses/{courseId}/announcements",
	//   "httpMethod": "GET",
	//   "id": "classroom.courses.announcements.list",
	//   "parameterOrder": [
	//     "courseId"
	//   ],
	//   "parameters": {
	//     "announcementStates": {
	//       "description": "Restriction on the `state` of announcements returned.\nIf this argument is left unspecified, the default value is `PUBLISHED`.",
	//       "enum": [
	//         "ANNOUNCEMENT_STATE_UNSPECIFIED",
	//         "PUBLISHED",
	//         "DRAFT",
	//         "DELETED"
	//       ],
	//       "location": "query",
	//       "repeated": true,
	//       "type": "string"
	//     },
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "orderBy": {
	//       "description": "Optional sort ordering for results. A comma-separated list of fields with\nan optional sort direction keyword. Supported field is `updateTime`.\nSupported direction keywords are `asc` and `desc`.\nIf not specified, `updateTime desc` is the default behavior.\nExamples: `updateTime asc`, `updateTime`",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum number of items to return. Zero or unspecified indicates that the\nserver may assign a maximum.\n\nThe server may return fewer than the specified number of results.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "nextPageToken\nvalue returned from a previous\nlist call,\nindicating that the subsequent page of results should be returned.\n\nThe list request\nmust be otherwise identical to the one that resulted in this token.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/announcements",
	//   "response": {
	//     "$ref": "ListAnnouncementsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.announcements",
	//     "https://www.googleapis.com/auth/classroom.announcements.readonly"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *CoursesAnnouncementsListCall) Pages(ctx context.Context, f func(*ListAnnouncementsResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "classroom.courses.announcements.modifyAssignees":

type CoursesAnnouncementsModifyAssigneesCall struct {
	s                                  *Service
	courseId                           string
	id                                 string
	modifyannouncementassigneesrequest *ModifyAnnouncementAssigneesRequest
	urlParams_                         gensupport.URLParams
	ctx_                               context.Context
	header_                            http.Header
}

// ModifyAssignees: Modifies assignee mode and options of an
// announcement.
//
// Only a teacher of the course that contains the announcement may
// call this method.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access the
// requested course or course work or for access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course or course work does not exist.
func (r *CoursesAnnouncementsService) ModifyAssignees(courseId string, id string, modifyannouncementassigneesrequest *ModifyAnnouncementAssigneesRequest) *CoursesAnnouncementsModifyAssigneesCall {
	c := &CoursesAnnouncementsModifyAssigneesCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.id = id
	c.modifyannouncementassigneesrequest = modifyannouncementassigneesrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesAnnouncementsModifyAssigneesCall) Fields(s ...googleapi.Field) *CoursesAnnouncementsModifyAssigneesCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesAnnouncementsModifyAssigneesCall) Context(ctx context.Context) *CoursesAnnouncementsModifyAssigneesCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesAnnouncementsModifyAssigneesCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesAnnouncementsModifyAssigneesCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.modifyannouncementassigneesrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/announcements/{id}:modifyAssignees")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
		"id":       c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.announcements.modifyAssignees" call.
// Exactly one of *Announcement or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Announcement.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *CoursesAnnouncementsModifyAssigneesCall) Do(opts ...googleapi.CallOption) (*Announcement, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Announcement{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Modifies assignee mode and options of an announcement.\n\nOnly a teacher of the course that contains the announcement may\ncall this method.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access the\nrequested course or course work or for access errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course or course work does not exist.",
	//   "flatPath": "v1/courses/{courseId}/announcements/{id}:modifyAssignees",
	//   "httpMethod": "POST",
	//   "id": "classroom.courses.announcements.modifyAssignees",
	//   "parameterOrder": [
	//     "courseId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "description": "Identifier of the announcement.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/announcements/{id}:modifyAssignees",
	//   "request": {
	//     "$ref": "ModifyAnnouncementAssigneesRequest"
	//   },
	//   "response": {
	//     "$ref": "Announcement"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.announcements"
	//   ]
	// }

}

// method id "classroom.courses.announcements.patch":

type CoursesAnnouncementsPatchCall struct {
	s            *Service
	courseId     string
	id           string
	announcement *Announcement
	urlParams_   gensupport.URLParams
	ctx_         context.Context
	header_      http.Header
}

// Patch: Updates one or more fields of an announcement.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting developer project did not
// create
// the corresponding announcement or for access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `FAILED_PRECONDITION` if the requested announcement has already
// been
// deleted.
// * `NOT_FOUND` if the requested course or announcement does not exist
func (r *CoursesAnnouncementsService) Patch(courseId string, id string, announcement *Announcement) *CoursesAnnouncementsPatchCall {
	c := &CoursesAnnouncementsPatchCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.id = id
	c.announcement = announcement
	return c
}

// UpdateMask sets the optional parameter "updateMask": Mask that
// identifies which fields on the announcement to update.
// This field is required to do an update. The update fails if
// invalid
// fields are specified. If a field supports empty values, it can be
// cleared
// by specifying it in the update mask and not in the Announcement
// object. If
// a field that does not support empty values is included in the update
// mask
// and not set in the Announcement object, an `INVALID_ARGUMENT` error
// will be
// returned.
//
// The following fields may be specified by teachers:
//
// * `text`
// * `state`
// * `scheduled_time`
func (c *CoursesAnnouncementsPatchCall) UpdateMask(updateMask string) *CoursesAnnouncementsPatchCall {
	c.urlParams_.Set("updateMask", updateMask)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesAnnouncementsPatchCall) Fields(s ...googleapi.Field) *CoursesAnnouncementsPatchCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesAnnouncementsPatchCall) Context(ctx context.Context) *CoursesAnnouncementsPatchCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesAnnouncementsPatchCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesAnnouncementsPatchCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.announcement)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/announcements/{id}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
		"id":       c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.announcements.patch" call.
// Exactly one of *Announcement or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Announcement.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *CoursesAnnouncementsPatchCall) Do(opts ...googleapi.CallOption) (*Announcement, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Announcement{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates one or more fields of an announcement.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting developer project did not create\nthe corresponding announcement or for access errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `FAILED_PRECONDITION` if the requested announcement has already been\ndeleted.\n* `NOT_FOUND` if the requested course or announcement does not exist",
	//   "flatPath": "v1/courses/{courseId}/announcements/{id}",
	//   "httpMethod": "PATCH",
	//   "id": "classroom.courses.announcements.patch",
	//   "parameterOrder": [
	//     "courseId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "description": "Identifier of the announcement.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "updateMask": {
	//       "description": "Mask that identifies which fields on the announcement to update.\nThis field is required to do an update. The update fails if invalid\nfields are specified. If a field supports empty values, it can be cleared\nby specifying it in the update mask and not in the Announcement object. If\na field that does not support empty values is included in the update mask\nand not set in the Announcement object, an `INVALID_ARGUMENT` error will be\nreturned.\n\nThe following fields may be specified by teachers:\n\n* `text`\n* `state`\n* `scheduled_time`",
	//       "format": "google-fieldmask",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/announcements/{id}",
	//   "request": {
	//     "$ref": "Announcement"
	//   },
	//   "response": {
	//     "$ref": "Announcement"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.announcements"
	//   ]
	// }

}

// method id "classroom.courses.courseWork.create":

type CoursesCourseWorkCreateCall struct {
	s          *Service
	courseId   string
	coursework *CourseWork
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Create: Creates course work.
//
// The resulting course work (and corresponding student submissions)
// are
// associated with the Developer Console project of the
// [OAuth client ID](https://support.google.com/cloud/answer/6158849)
// used to
// make the request. Classroom API requests to modify course work and
// student
// submissions must be made with an OAuth client ID from the
// associated
// Developer Console project.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access the
// requested course, create course work in the requested course, share
// a
// Drive attachment, or for access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course does not exist.
// * `FAILED_PRECONDITION` for the following request error:
//     * AttachmentNotVisible
func (r *CoursesCourseWorkService) Create(courseId string, coursework *CourseWork) *CoursesCourseWorkCreateCall {
	c := &CoursesCourseWorkCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.coursework = coursework
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesCourseWorkCreateCall) Fields(s ...googleapi.Field) *CoursesCourseWorkCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesCourseWorkCreateCall) Context(ctx context.Context) *CoursesCourseWorkCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesCourseWorkCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesCourseWorkCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.coursework)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/courseWork")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.courseWork.create" call.
// Exactly one of *CourseWork or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *CourseWork.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *CoursesCourseWorkCreateCall) Do(opts ...googleapi.CallOption) (*CourseWork, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &CourseWork{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Creates course work.\n\nThe resulting course work (and corresponding student submissions) are\nassociated with the Developer Console project of the\n[OAuth client ID](https://support.google.com/cloud/answer/6158849) used to\nmake the request. Classroom API requests to modify course work and student\nsubmissions must be made with an OAuth client ID from the associated\nDeveloper Console project.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access the\nrequested course, create course work in the requested course, share a\nDrive attachment, or for access errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course does not exist.\n* `FAILED_PRECONDITION` for the following request error:\n    * AttachmentNotVisible",
	//   "flatPath": "v1/courses/{courseId}/courseWork",
	//   "httpMethod": "POST",
	//   "id": "classroom.courses.courseWork.create",
	//   "parameterOrder": [
	//     "courseId"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/courseWork",
	//   "request": {
	//     "$ref": "CourseWork"
	//   },
	//   "response": {
	//     "$ref": "CourseWork"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.coursework.students"
	//   ]
	// }

}

// method id "classroom.courses.courseWork.delete":

type CoursesCourseWorkDeleteCall struct {
	s          *Service
	courseId   string
	id         string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes a course work.
//
// This request must be made by the Developer Console project of
// the
// [OAuth client ID](https://support.google.com/cloud/answer/6158849)
// used to
// create the corresponding course work item.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting developer project did not
// create
// the corresponding course work, if the requesting user is not
// permitted
// to delete the requested course or for access errors.
// * `FAILED_PRECONDITION` if the requested course work has already
// been
// deleted.
// * `NOT_FOUND` if no course exists with the requested ID.
func (r *CoursesCourseWorkService) Delete(courseId string, id string) *CoursesCourseWorkDeleteCall {
	c := &CoursesCourseWorkDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.id = id
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesCourseWorkDeleteCall) Fields(s ...googleapi.Field) *CoursesCourseWorkDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesCourseWorkDeleteCall) Context(ctx context.Context) *CoursesCourseWorkDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesCourseWorkDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesCourseWorkDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/courseWork/{id}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
		"id":       c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.courseWork.delete" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesCourseWorkDeleteCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Deletes a course work.\n\nThis request must be made by the Developer Console project of the\n[OAuth client ID](https://support.google.com/cloud/answer/6158849) used to\ncreate the corresponding course work item.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting developer project did not create\nthe corresponding course work, if the requesting user is not permitted\nto delete the requested course or for access errors.\n* `FAILED_PRECONDITION` if the requested course work has already been\ndeleted.\n* `NOT_FOUND` if no course exists with the requested ID.",
	//   "flatPath": "v1/courses/{courseId}/courseWork/{id}",
	//   "httpMethod": "DELETE",
	//   "id": "classroom.courses.courseWork.delete",
	//   "parameterOrder": [
	//     "courseId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "description": "Identifier of the course work to delete.\nThis identifier is a Classroom-assigned identifier.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/courseWork/{id}",
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.coursework.students"
	//   ]
	// }

}

// method id "classroom.courses.courseWork.get":

type CoursesCourseWorkGetCall struct {
	s            *Service
	courseId     string
	id           string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Returns course work.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access the
// requested course or course work, or for access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course or course work does not exist.
func (r *CoursesCourseWorkService) Get(courseId string, id string) *CoursesCourseWorkGetCall {
	c := &CoursesCourseWorkGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.id = id
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesCourseWorkGetCall) Fields(s ...googleapi.Field) *CoursesCourseWorkGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *CoursesCourseWorkGetCall) IfNoneMatch(entityTag string) *CoursesCourseWorkGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesCourseWorkGetCall) Context(ctx context.Context) *CoursesCourseWorkGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesCourseWorkGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesCourseWorkGetCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/courseWork/{id}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
		"id":       c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.courseWork.get" call.
// Exactly one of *CourseWork or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *CourseWork.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *CoursesCourseWorkGetCall) Do(opts ...googleapi.CallOption) (*CourseWork, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &CourseWork{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns course work.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access the\nrequested course or course work, or for access errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course or course work does not exist.",
	//   "flatPath": "v1/courses/{courseId}/courseWork/{id}",
	//   "httpMethod": "GET",
	//   "id": "classroom.courses.courseWork.get",
	//   "parameterOrder": [
	//     "courseId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "description": "Identifier of the course work.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/courseWork/{id}",
	//   "response": {
	//     "$ref": "CourseWork"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.coursework.me",
	//     "https://www.googleapis.com/auth/classroom.coursework.me.readonly",
	//     "https://www.googleapis.com/auth/classroom.coursework.students",
	//     "https://www.googleapis.com/auth/classroom.coursework.students.readonly"
	//   ]
	// }

}

// method id "classroom.courses.courseWork.list":

type CoursesCourseWorkListCall struct {
	s            *Service
	courseId     string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Returns a list of course work that the requester is permitted
// to view.
//
// Course students may only view `PUBLISHED` course work. Course
// teachers
// and domain administrators may view all course work.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access
// the requested course or for access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course does not exist.
func (r *CoursesCourseWorkService) List(courseId string) *CoursesCourseWorkListCall {
	c := &CoursesCourseWorkListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	return c
}

// CourseWorkStates sets the optional parameter "courseWorkStates":
// Restriction on the work status to return. Only courseWork that
// matches
// is returned. If unspecified, items with a work status of
// `PUBLISHED`
// is returned.
//
// Possible values:
//   "COURSE_WORK_STATE_UNSPECIFIED"
//   "PUBLISHED"
//   "DRAFT"
//   "DELETED"
func (c *CoursesCourseWorkListCall) CourseWorkStates(courseWorkStates ...string) *CoursesCourseWorkListCall {
	c.urlParams_.SetMulti("courseWorkStates", append([]string{}, courseWorkStates...))
	return c
}

// OrderBy sets the optional parameter "orderBy": Optional sort ordering
// for results. A comma-separated list of fields with
// an optional sort direction keyword. Supported fields are
// `updateTime`
// and `dueDate`. Supported direction keywords are `asc` and `desc`.
// If not specified, `updateTime desc` is the default
// behavior.
// Examples: `dueDate asc,updateTime desc`, `updateTime,dueDate desc`
func (c *CoursesCourseWorkListCall) OrderBy(orderBy string) *CoursesCourseWorkListCall {
	c.urlParams_.Set("orderBy", orderBy)
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum number of
// items to return. Zero or unspecified indicates that the
// server may assign a maximum.
//
// The server may return fewer than the specified number of results.
func (c *CoursesCourseWorkListCall) PageSize(pageSize int64) *CoursesCourseWorkListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken":
// nextPageToken
// value returned from a previous
// list call,
// indicating that the subsequent page of results should be
// returned.
//
// The list request
// must be otherwise identical to the one that resulted in this token.
func (c *CoursesCourseWorkListCall) PageToken(pageToken string) *CoursesCourseWorkListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesCourseWorkListCall) Fields(s ...googleapi.Field) *CoursesCourseWorkListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *CoursesCourseWorkListCall) IfNoneMatch(entityTag string) *CoursesCourseWorkListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesCourseWorkListCall) Context(ctx context.Context) *CoursesCourseWorkListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesCourseWorkListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesCourseWorkListCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/courseWork")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.courseWork.list" call.
// Exactly one of *ListCourseWorkResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListCourseWorkResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *CoursesCourseWorkListCall) Do(opts ...googleapi.CallOption) (*ListCourseWorkResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListCourseWorkResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a list of course work that the requester is permitted to view.\n\nCourse students may only view `PUBLISHED` course work. Course teachers\nand domain administrators may view all course work.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access\nthe requested course or for access errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course does not exist.",
	//   "flatPath": "v1/courses/{courseId}/courseWork",
	//   "httpMethod": "GET",
	//   "id": "classroom.courses.courseWork.list",
	//   "parameterOrder": [
	//     "courseId"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "courseWorkStates": {
	//       "description": "Restriction on the work status to return. Only courseWork that matches\nis returned. If unspecified, items with a work status of `PUBLISHED`\nis returned.",
	//       "enum": [
	//         "COURSE_WORK_STATE_UNSPECIFIED",
	//         "PUBLISHED",
	//         "DRAFT",
	//         "DELETED"
	//       ],
	//       "location": "query",
	//       "repeated": true,
	//       "type": "string"
	//     },
	//     "orderBy": {
	//       "description": "Optional sort ordering for results. A comma-separated list of fields with\nan optional sort direction keyword. Supported fields are `updateTime`\nand `dueDate`. Supported direction keywords are `asc` and `desc`.\nIf not specified, `updateTime desc` is the default behavior.\nExamples: `dueDate asc,updateTime desc`, `updateTime,dueDate desc`",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum number of items to return. Zero or unspecified indicates that the\nserver may assign a maximum.\n\nThe server may return fewer than the specified number of results.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "nextPageToken\nvalue returned from a previous\nlist call,\nindicating that the subsequent page of results should be returned.\n\nThe list request\nmust be otherwise identical to the one that resulted in this token.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/courseWork",
	//   "response": {
	//     "$ref": "ListCourseWorkResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.coursework.me",
	//     "https://www.googleapis.com/auth/classroom.coursework.me.readonly",
	//     "https://www.googleapis.com/auth/classroom.coursework.students",
	//     "https://www.googleapis.com/auth/classroom.coursework.students.readonly"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *CoursesCourseWorkListCall) Pages(ctx context.Context, f func(*ListCourseWorkResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "classroom.courses.courseWork.modifyAssignees":

type CoursesCourseWorkModifyAssigneesCall struct {
	s                                *Service
	courseId                         string
	id                               string
	modifycourseworkassigneesrequest *ModifyCourseWorkAssigneesRequest
	urlParams_                       gensupport.URLParams
	ctx_                             context.Context
	header_                          http.Header
}

// ModifyAssignees: Modifies assignee mode and options of a
// coursework.
//
// Only a teacher of the course that contains the coursework may
// call this method.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access the
// requested course or course work or for access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course or course work does not exist.
func (r *CoursesCourseWorkService) ModifyAssignees(courseId string, id string, modifycourseworkassigneesrequest *ModifyCourseWorkAssigneesRequest) *CoursesCourseWorkModifyAssigneesCall {
	c := &CoursesCourseWorkModifyAssigneesCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.id = id
	c.modifycourseworkassigneesrequest = modifycourseworkassigneesrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesCourseWorkModifyAssigneesCall) Fields(s ...googleapi.Field) *CoursesCourseWorkModifyAssigneesCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesCourseWorkModifyAssigneesCall) Context(ctx context.Context) *CoursesCourseWorkModifyAssigneesCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesCourseWorkModifyAssigneesCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesCourseWorkModifyAssigneesCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.modifycourseworkassigneesrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/courseWork/{id}:modifyAssignees")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
		"id":       c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.courseWork.modifyAssignees" call.
// Exactly one of *CourseWork or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *CourseWork.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *CoursesCourseWorkModifyAssigneesCall) Do(opts ...googleapi.CallOption) (*CourseWork, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &CourseWork{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Modifies assignee mode and options of a coursework.\n\nOnly a teacher of the course that contains the coursework may\ncall this method.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access the\nrequested course or course work or for access errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course or course work does not exist.",
	//   "flatPath": "v1/courses/{courseId}/courseWork/{id}:modifyAssignees",
	//   "httpMethod": "POST",
	//   "id": "classroom.courses.courseWork.modifyAssignees",
	//   "parameterOrder": [
	//     "courseId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "description": "Identifier of the coursework.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/courseWork/{id}:modifyAssignees",
	//   "request": {
	//     "$ref": "ModifyCourseWorkAssigneesRequest"
	//   },
	//   "response": {
	//     "$ref": "CourseWork"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.coursework.students"
	//   ]
	// }

}

// method id "classroom.courses.courseWork.patch":

type CoursesCourseWorkPatchCall struct {
	s          *Service
	courseId   string
	id         string
	coursework *CourseWork
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Patch: Updates one or more fields of a course work.
//
// See google.classroom.v1.CourseWork for details
// of which fields may be updated and who may change them.
//
// This request must be made by the Developer Console project of
// the
// [OAuth client ID](https://support.google.com/cloud/answer/6158849)
// used to
// create the corresponding course work item.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting developer project did not
// create
// the corresponding course work, if the user is not permitted to make
// the
// requested modification to the student submission, or for
// access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `FAILED_PRECONDITION` if the requested course work has already
// been
// deleted.
// * `NOT_FOUND` if the requested course, course work, or student
// submission
// does not exist.
func (r *CoursesCourseWorkService) Patch(courseId string, id string, coursework *CourseWork) *CoursesCourseWorkPatchCall {
	c := &CoursesCourseWorkPatchCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.id = id
	c.coursework = coursework
	return c
}

// UpdateMask sets the optional parameter "updateMask": Mask that
// identifies which fields on the course work to update.
// This field is required to do an update. The update fails if
// invalid
// fields are specified. If a field supports empty values, it can be
// cleared
// by specifying it in the update mask and not in the CourseWork object.
// If a
// field that does not support empty values is included in the update
// mask and
// not set in the CourseWork object, an `INVALID_ARGUMENT` error will
// be
// returned.
//
// The following fields may be specified by teachers:
//
// * `title`
// * `description`
// * `state`
// * `due_date`
// * `due_time`
// * `max_points`
// * `scheduled_time`
// * `submission_modification_mode`
func (c *CoursesCourseWorkPatchCall) UpdateMask(updateMask string) *CoursesCourseWorkPatchCall {
	c.urlParams_.Set("updateMask", updateMask)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesCourseWorkPatchCall) Fields(s ...googleapi.Field) *CoursesCourseWorkPatchCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesCourseWorkPatchCall) Context(ctx context.Context) *CoursesCourseWorkPatchCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesCourseWorkPatchCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesCourseWorkPatchCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.coursework)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/courseWork/{id}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
		"id":       c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.courseWork.patch" call.
// Exactly one of *CourseWork or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *CourseWork.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *CoursesCourseWorkPatchCall) Do(opts ...googleapi.CallOption) (*CourseWork, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &CourseWork{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates one or more fields of a course work.\n\nSee google.classroom.v1.CourseWork for details\nof which fields may be updated and who may change them.\n\nThis request must be made by the Developer Console project of the\n[OAuth client ID](https://support.google.com/cloud/answer/6158849) used to\ncreate the corresponding course work item.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting developer project did not create\nthe corresponding course work, if the user is not permitted to make the\nrequested modification to the student submission, or for\naccess errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `FAILED_PRECONDITION` if the requested course work has already been\ndeleted.\n* `NOT_FOUND` if the requested course, course work, or student submission\ndoes not exist.",
	//   "flatPath": "v1/courses/{courseId}/courseWork/{id}",
	//   "httpMethod": "PATCH",
	//   "id": "classroom.courses.courseWork.patch",
	//   "parameterOrder": [
	//     "courseId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "description": "Identifier of the course work.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "updateMask": {
	//       "description": "Mask that identifies which fields on the course work to update.\nThis field is required to do an update. The update fails if invalid\nfields are specified. If a field supports empty values, it can be cleared\nby specifying it in the update mask and not in the CourseWork object. If a\nfield that does not support empty values is included in the update mask and\nnot set in the CourseWork object, an `INVALID_ARGUMENT` error will be\nreturned.\n\nThe following fields may be specified by teachers:\n\n* `title`\n* `description`\n* `state`\n* `due_date`\n* `due_time`\n* `max_points`\n* `scheduled_time`\n* `submission_modification_mode`",
	//       "format": "google-fieldmask",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/courseWork/{id}",
	//   "request": {
	//     "$ref": "CourseWork"
	//   },
	//   "response": {
	//     "$ref": "CourseWork"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.coursework.students"
	//   ]
	// }

}

// method id "classroom.courses.courseWork.studentSubmissions.get":

type CoursesCourseWorkStudentSubmissionsGetCall struct {
	s            *Service
	courseId     string
	courseWorkId string
	id           string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Returns a student submission.
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access the
// requested course, course work, or student submission or for
// access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course, course work, or student
// submission
// does not exist.
func (r *CoursesCourseWorkStudentSubmissionsService) Get(courseId string, courseWorkId string, id string) *CoursesCourseWorkStudentSubmissionsGetCall {
	c := &CoursesCourseWorkStudentSubmissionsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.courseWorkId = courseWorkId
	c.id = id
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesCourseWorkStudentSubmissionsGetCall) Fields(s ...googleapi.Field) *CoursesCourseWorkStudentSubmissionsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *CoursesCourseWorkStudentSubmissionsGetCall) IfNoneMatch(entityTag string) *CoursesCourseWorkStudentSubmissionsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesCourseWorkStudentSubmissionsGetCall) Context(ctx context.Context) *CoursesCourseWorkStudentSubmissionsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesCourseWorkStudentSubmissionsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesCourseWorkStudentSubmissionsGetCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId":     c.courseId,
		"courseWorkId": c.courseWorkId,
		"id":           c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.courseWork.studentSubmissions.get" call.
// Exactly one of *StudentSubmission or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *StudentSubmission.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *CoursesCourseWorkStudentSubmissionsGetCall) Do(opts ...googleapi.CallOption) (*StudentSubmission, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &StudentSubmission{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a student submission.\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access the\nrequested course, course work, or student submission or for\naccess errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course, course work, or student submission\ndoes not exist.",
	//   "flatPath": "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}",
	//   "httpMethod": "GET",
	//   "id": "classroom.courses.courseWork.studentSubmissions.get",
	//   "parameterOrder": [
	//     "courseId",
	//     "courseWorkId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "courseWorkId": {
	//       "description": "Identifier of the course work.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "description": "Identifier of the student submission.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}",
	//   "response": {
	//     "$ref": "StudentSubmission"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.coursework.me",
	//     "https://www.googleapis.com/auth/classroom.coursework.me.readonly",
	//     "https://www.googleapis.com/auth/classroom.coursework.students",
	//     "https://www.googleapis.com/auth/classroom.coursework.students.readonly",
	//     "https://www.googleapis.com/auth/classroom.student-submissions.me.readonly",
	//     "https://www.googleapis.com/auth/classroom.student-submissions.students.readonly"
	//   ]
	// }

}

// method id "classroom.courses.courseWork.studentSubmissions.list":

type CoursesCourseWorkStudentSubmissionsListCall struct {
	s            *Service
	courseId     string
	courseWorkId string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Returns a list of student submissions that the requester is
// permitted to
// view, factoring in the OAuth scopes of the request.
// `-` may be specified as the `course_work_id` to include
// student
// submissions for multiple course work items.
//
// Course students may only view their own work. Course teachers
// and domain administrators may view all student submissions.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access the
// requested course or course work, or for access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course does not exist.
func (r *CoursesCourseWorkStudentSubmissionsService) List(courseId string, courseWorkId string) *CoursesCourseWorkStudentSubmissionsListCall {
	c := &CoursesCourseWorkStudentSubmissionsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.courseWorkId = courseWorkId
	return c
}

// Late sets the optional parameter "late": Requested lateness value. If
// specified, returned student submissions are
// restricted by the requested value.
// If unspecified, submissions are returned regardless of `late` value.
//
// Possible values:
//   "LATE_VALUES_UNSPECIFIED"
//   "LATE_ONLY"
//   "NOT_LATE_ONLY"
func (c *CoursesCourseWorkStudentSubmissionsListCall) Late(late string) *CoursesCourseWorkStudentSubmissionsListCall {
	c.urlParams_.Set("late", late)
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum number of
// items to return. Zero or unspecified indicates that the
// server may assign a maximum.
//
// The server may return fewer than the specified number of results.
func (c *CoursesCourseWorkStudentSubmissionsListCall) PageSize(pageSize int64) *CoursesCourseWorkStudentSubmissionsListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken":
// nextPageToken
// value returned from a previous
// list call,
// indicating that the subsequent page of results should be
// returned.
//
// The list request
// must be otherwise identical to the one that resulted in this token.
func (c *CoursesCourseWorkStudentSubmissionsListCall) PageToken(pageToken string) *CoursesCourseWorkStudentSubmissionsListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// States sets the optional parameter "states": Requested submission
// states. If specified, returned student submissions
// match one of the specified submission states.
//
// Possible values:
//   "SUBMISSION_STATE_UNSPECIFIED"
//   "NEW"
//   "CREATED"
//   "TURNED_IN"
//   "RETURNED"
//   "RECLAIMED_BY_STUDENT"
func (c *CoursesCourseWorkStudentSubmissionsListCall) States(states ...string) *CoursesCourseWorkStudentSubmissionsListCall {
	c.urlParams_.SetMulti("states", append([]string{}, states...))
	return c
}

// UserId sets the optional parameter "userId": Optional argument to
// restrict returned student work to those owned by the
// student with the specified identifier. The identifier can be one of
// the
// following:
//
// * the numeric identifier for the user
// * the email address of the user
// * the string literal "me", indicating the requesting user
func (c *CoursesCourseWorkStudentSubmissionsListCall) UserId(userId string) *CoursesCourseWorkStudentSubmissionsListCall {
	c.urlParams_.Set("userId", userId)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesCourseWorkStudentSubmissionsListCall) Fields(s ...googleapi.Field) *CoursesCourseWorkStudentSubmissionsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *CoursesCourseWorkStudentSubmissionsListCall) IfNoneMatch(entityTag string) *CoursesCourseWorkStudentSubmissionsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesCourseWorkStudentSubmissionsListCall) Context(ctx context.Context) *CoursesCourseWorkStudentSubmissionsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesCourseWorkStudentSubmissionsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesCourseWorkStudentSubmissionsListCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId":     c.courseId,
		"courseWorkId": c.courseWorkId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.courseWork.studentSubmissions.list" call.
// Exactly one of *ListStudentSubmissionsResponse or error will be
// non-nil. Any non-2xx status code is an error. Response headers are in
// either *ListStudentSubmissionsResponse.ServerResponse.Header or (if a
// response was returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *CoursesCourseWorkStudentSubmissionsListCall) Do(opts ...googleapi.CallOption) (*ListStudentSubmissionsResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListStudentSubmissionsResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a list of student submissions that the requester is permitted to\nview, factoring in the OAuth scopes of the request.\n`-` may be specified as the `course_work_id` to include student\nsubmissions for multiple course work items.\n\nCourse students may only view their own work. Course teachers\nand domain administrators may view all student submissions.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access the\nrequested course or course work, or for access errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course does not exist.",
	//   "flatPath": "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions",
	//   "httpMethod": "GET",
	//   "id": "classroom.courses.courseWork.studentSubmissions.list",
	//   "parameterOrder": [
	//     "courseId",
	//     "courseWorkId"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "courseWorkId": {
	//       "description": "Identifier of the student work to request.\nThis may be set to the string literal `\"-\"` to request student work for\nall course work in the specified course.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "late": {
	//       "description": "Requested lateness value. If specified, returned student submissions are\nrestricted by the requested value.\nIf unspecified, submissions are returned regardless of `late` value.",
	//       "enum": [
	//         "LATE_VALUES_UNSPECIFIED",
	//         "LATE_ONLY",
	//         "NOT_LATE_ONLY"
	//       ],
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum number of items to return. Zero or unspecified indicates that the\nserver may assign a maximum.\n\nThe server may return fewer than the specified number of results.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "nextPageToken\nvalue returned from a previous\nlist call,\nindicating that the subsequent page of results should be returned.\n\nThe list request\nmust be otherwise identical to the one that resulted in this token.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "states": {
	//       "description": "Requested submission states. If specified, returned student submissions\nmatch one of the specified submission states.",
	//       "enum": [
	//         "SUBMISSION_STATE_UNSPECIFIED",
	//         "NEW",
	//         "CREATED",
	//         "TURNED_IN",
	//         "RETURNED",
	//         "RECLAIMED_BY_STUDENT"
	//       ],
	//       "location": "query",
	//       "repeated": true,
	//       "type": "string"
	//     },
	//     "userId": {
	//       "description": "Optional argument to restrict returned student work to those owned by the\nstudent with the specified identifier. The identifier can be one of the\nfollowing:\n\n* the numeric identifier for the user\n* the email address of the user\n* the string literal `\"me\"`, indicating the requesting user",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions",
	//   "response": {
	//     "$ref": "ListStudentSubmissionsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.coursework.me",
	//     "https://www.googleapis.com/auth/classroom.coursework.me.readonly",
	//     "https://www.googleapis.com/auth/classroom.coursework.students",
	//     "https://www.googleapis.com/auth/classroom.coursework.students.readonly",
	//     "https://www.googleapis.com/auth/classroom.student-submissions.me.readonly",
	//     "https://www.googleapis.com/auth/classroom.student-submissions.students.readonly"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *CoursesCourseWorkStudentSubmissionsListCall) Pages(ctx context.Context, f func(*ListStudentSubmissionsResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "classroom.courses.courseWork.studentSubmissions.modifyAttachments":

type CoursesCourseWorkStudentSubmissionsModifyAttachmentsCall struct {
	s                        *Service
	courseId                 string
	courseWorkId             string
	id                       string
	modifyattachmentsrequest *ModifyAttachmentsRequest
	urlParams_               gensupport.URLParams
	ctx_                     context.Context
	header_                  http.Header
}

// ModifyAttachments: Modifies attachments of student
// submission.
//
// Attachments may only be added to student submissions belonging to
// course
// work objects with a `workType` of `ASSIGNMENT`.
//
// This request must be made by the Developer Console project of
// the
// [OAuth client ID](https://support.google.com/cloud/answer/6158849)
// used to
// create the corresponding course work item.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access the
// requested course or course work, if the user is not permitted to
// modify
// attachments on the requested student submission, or for
// access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course, course work, or student
// submission
// does not exist.
func (r *CoursesCourseWorkStudentSubmissionsService) ModifyAttachments(courseId string, courseWorkId string, id string, modifyattachmentsrequest *ModifyAttachmentsRequest) *CoursesCourseWorkStudentSubmissionsModifyAttachmentsCall {
	c := &CoursesCourseWorkStudentSubmissionsModifyAttachmentsCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.courseWorkId = courseWorkId
	c.id = id
	c.modifyattachmentsrequest = modifyattachmentsrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesCourseWorkStudentSubmissionsModifyAttachmentsCall) Fields(s ...googleapi.Field) *CoursesCourseWorkStudentSubmissionsModifyAttachmentsCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesCourseWorkStudentSubmissionsModifyAttachmentsCall) Context(ctx context.Context) *CoursesCourseWorkStudentSubmissionsModifyAttachmentsCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesCourseWorkStudentSubmissionsModifyAttachmentsCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesCourseWorkStudentSubmissionsModifyAttachmentsCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.modifyattachmentsrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}:modifyAttachments")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId":     c.courseId,
		"courseWorkId": c.courseWorkId,
		"id":           c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.courseWork.studentSubmissions.modifyAttachments" call.
// Exactly one of *StudentSubmission or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *StudentSubmission.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *CoursesCourseWorkStudentSubmissionsModifyAttachmentsCall) Do(opts ...googleapi.CallOption) (*StudentSubmission, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &StudentSubmission{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Modifies attachments of student submission.\n\nAttachments may only be added to student submissions belonging to course\nwork objects with a `workType` of `ASSIGNMENT`.\n\nThis request must be made by the Developer Console project of the\n[OAuth client ID](https://support.google.com/cloud/answer/6158849) used to\ncreate the corresponding course work item.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access the\nrequested course or course work, if the user is not permitted to modify\nattachments on the requested student submission, or for\naccess errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course, course work, or student submission\ndoes not exist.",
	//   "flatPath": "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}:modifyAttachments",
	//   "httpMethod": "POST",
	//   "id": "classroom.courses.courseWork.studentSubmissions.modifyAttachments",
	//   "parameterOrder": [
	//     "courseId",
	//     "courseWorkId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "courseWorkId": {
	//       "description": "Identifier of the course work.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "description": "Identifier of the student submission.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}:modifyAttachments",
	//   "request": {
	//     "$ref": "ModifyAttachmentsRequest"
	//   },
	//   "response": {
	//     "$ref": "StudentSubmission"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.coursework.me",
	//     "https://www.googleapis.com/auth/classroom.coursework.students"
	//   ]
	// }

}

// method id "classroom.courses.courseWork.studentSubmissions.patch":

type CoursesCourseWorkStudentSubmissionsPatchCall struct {
	s                 *Service
	courseId          string
	courseWorkId      string
	id                string
	studentsubmission *StudentSubmission
	urlParams_        gensupport.URLParams
	ctx_              context.Context
	header_           http.Header
}

// Patch: Updates one or more fields of a student submission.
//
// See google.classroom.v1.StudentSubmission for details
// of which fields may be updated and who may change them.
//
// This request must be made by the Developer Console project of
// the
// [OAuth client ID](https://support.google.com/cloud/answer/6158849)
// used to
// create the corresponding course work item.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting developer project did not
// create
// the corresponding course work, if the user is not permitted to make
// the
// requested modification to the student submission, or for
// access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course, course work, or student
// submission
// does not exist.
func (r *CoursesCourseWorkStudentSubmissionsService) Patch(courseId string, courseWorkId string, id string, studentsubmission *StudentSubmission) *CoursesCourseWorkStudentSubmissionsPatchCall {
	c := &CoursesCourseWorkStudentSubmissionsPatchCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.courseWorkId = courseWorkId
	c.id = id
	c.studentsubmission = studentsubmission
	return c
}

// UpdateMask sets the optional parameter "updateMask": Mask that
// identifies which fields on the student submission to update.
// This field is required to do an update. The update fails if
// invalid
// fields are specified.
//
// The following fields may be specified by teachers:
//
// * `draft_grade`
// * `assigned_grade`
func (c *CoursesCourseWorkStudentSubmissionsPatchCall) UpdateMask(updateMask string) *CoursesCourseWorkStudentSubmissionsPatchCall {
	c.urlParams_.Set("updateMask", updateMask)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesCourseWorkStudentSubmissionsPatchCall) Fields(s ...googleapi.Field) *CoursesCourseWorkStudentSubmissionsPatchCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesCourseWorkStudentSubmissionsPatchCall) Context(ctx context.Context) *CoursesCourseWorkStudentSubmissionsPatchCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesCourseWorkStudentSubmissionsPatchCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesCourseWorkStudentSubmissionsPatchCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.studentsubmission)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId":     c.courseId,
		"courseWorkId": c.courseWorkId,
		"id":           c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.courseWork.studentSubmissions.patch" call.
// Exactly one of *StudentSubmission or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *StudentSubmission.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *CoursesCourseWorkStudentSubmissionsPatchCall) Do(opts ...googleapi.CallOption) (*StudentSubmission, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &StudentSubmission{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates one or more fields of a student submission.\n\nSee google.classroom.v1.StudentSubmission for details\nof which fields may be updated and who may change them.\n\nThis request must be made by the Developer Console project of the\n[OAuth client ID](https://support.google.com/cloud/answer/6158849) used to\ncreate the corresponding course work item.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting developer project did not create\nthe corresponding course work, if the user is not permitted to make the\nrequested modification to the student submission, or for\naccess errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course, course work, or student submission\ndoes not exist.",
	//   "flatPath": "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}",
	//   "httpMethod": "PATCH",
	//   "id": "classroom.courses.courseWork.studentSubmissions.patch",
	//   "parameterOrder": [
	//     "courseId",
	//     "courseWorkId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "courseWorkId": {
	//       "description": "Identifier of the course work.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "description": "Identifier of the student submission.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "updateMask": {
	//       "description": "Mask that identifies which fields on the student submission to update.\nThis field is required to do an update. The update fails if invalid\nfields are specified.\n\nThe following fields may be specified by teachers:\n\n* `draft_grade`\n* `assigned_grade`",
	//       "format": "google-fieldmask",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}",
	//   "request": {
	//     "$ref": "StudentSubmission"
	//   },
	//   "response": {
	//     "$ref": "StudentSubmission"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.coursework.me",
	//     "https://www.googleapis.com/auth/classroom.coursework.students"
	//   ]
	// }

}

// method id "classroom.courses.courseWork.studentSubmissions.reclaim":

type CoursesCourseWorkStudentSubmissionsReclaimCall struct {
	s                               *Service
	courseId                        string
	courseWorkId                    string
	id                              string
	reclaimstudentsubmissionrequest *ReclaimStudentSubmissionRequest
	urlParams_                      gensupport.URLParams
	ctx_                            context.Context
	header_                         http.Header
}

// Reclaim: Reclaims a student submission on behalf of the student that
// owns it.
//
// Reclaiming a student submission transfers ownership of attached
// Drive
// files to the student and update the submission state.
//
// Only the student that owns the requested student submission may call
// this
// method, and only for a student submission that has been turned
// in.
//
// This request must be made by the Developer Console project of
// the
// [OAuth client ID](https://support.google.com/cloud/answer/6158849)
// used to
// create the corresponding course work item.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access the
// requested course or course work, unsubmit the requested student
// submission,
// or for access errors.
// * `FAILED_PRECONDITION` if the student submission has not been turned
// in.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course, course work, or student
// submission
// does not exist.
func (r *CoursesCourseWorkStudentSubmissionsService) Reclaim(courseId string, courseWorkId string, id string, reclaimstudentsubmissionrequest *ReclaimStudentSubmissionRequest) *CoursesCourseWorkStudentSubmissionsReclaimCall {
	c := &CoursesCourseWorkStudentSubmissionsReclaimCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.courseWorkId = courseWorkId
	c.id = id
	c.reclaimstudentsubmissionrequest = reclaimstudentsubmissionrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesCourseWorkStudentSubmissionsReclaimCall) Fields(s ...googleapi.Field) *CoursesCourseWorkStudentSubmissionsReclaimCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesCourseWorkStudentSubmissionsReclaimCall) Context(ctx context.Context) *CoursesCourseWorkStudentSubmissionsReclaimCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesCourseWorkStudentSubmissionsReclaimCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesCourseWorkStudentSubmissionsReclaimCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.reclaimstudentsubmissionrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}:reclaim")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId":     c.courseId,
		"courseWorkId": c.courseWorkId,
		"id":           c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.courseWork.studentSubmissions.reclaim" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesCourseWorkStudentSubmissionsReclaimCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Reclaims a student submission on behalf of the student that owns it.\n\nReclaiming a student submission transfers ownership of attached Drive\nfiles to the student and update the submission state.\n\nOnly the student that owns the requested student submission may call this\nmethod, and only for a student submission that has been turned in.\n\nThis request must be made by the Developer Console project of the\n[OAuth client ID](https://support.google.com/cloud/answer/6158849) used to\ncreate the corresponding course work item.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access the\nrequested course or course work, unsubmit the requested student submission,\nor for access errors.\n* `FAILED_PRECONDITION` if the student submission has not been turned in.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course, course work, or student submission\ndoes not exist.",
	//   "flatPath": "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}:reclaim",
	//   "httpMethod": "POST",
	//   "id": "classroom.courses.courseWork.studentSubmissions.reclaim",
	//   "parameterOrder": [
	//     "courseId",
	//     "courseWorkId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "courseWorkId": {
	//       "description": "Identifier of the course work.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "description": "Identifier of the student submission.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}:reclaim",
	//   "request": {
	//     "$ref": "ReclaimStudentSubmissionRequest"
	//   },
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.coursework.me"
	//   ]
	// }

}

// method id "classroom.courses.courseWork.studentSubmissions.return":

type CoursesCourseWorkStudentSubmissionsReturnCall struct {
	s                              *Service
	courseId                       string
	courseWorkId                   string
	id                             string
	returnstudentsubmissionrequest *ReturnStudentSubmissionRequest
	urlParams_                     gensupport.URLParams
	ctx_                           context.Context
	header_                        http.Header
}

// Return: Returns a student submission.
//
// Returning a student submission transfers ownership of attached
// Drive
// files to the student and may also update the submission state.
// Unlike the Classroom application, returning a student submission does
// not
// set assignedGrade to the draftGrade value.
//
// Only a teacher of the course that contains the requested student
// submission
// may call this method.
//
// This request must be made by the Developer Console project of
// the
// [OAuth client ID](https://support.google.com/cloud/answer/6158849)
// used to
// create the corresponding course work item.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access the
// requested course or course work, return the requested student
// submission,
// or for access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course, course work, or student
// submission
// does not exist.
func (r *CoursesCourseWorkStudentSubmissionsService) Return(courseId string, courseWorkId string, id string, returnstudentsubmissionrequest *ReturnStudentSubmissionRequest) *CoursesCourseWorkStudentSubmissionsReturnCall {
	c := &CoursesCourseWorkStudentSubmissionsReturnCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.courseWorkId = courseWorkId
	c.id = id
	c.returnstudentsubmissionrequest = returnstudentsubmissionrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesCourseWorkStudentSubmissionsReturnCall) Fields(s ...googleapi.Field) *CoursesCourseWorkStudentSubmissionsReturnCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesCourseWorkStudentSubmissionsReturnCall) Context(ctx context.Context) *CoursesCourseWorkStudentSubmissionsReturnCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesCourseWorkStudentSubmissionsReturnCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesCourseWorkStudentSubmissionsReturnCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.returnstudentsubmissionrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}:return")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId":     c.courseId,
		"courseWorkId": c.courseWorkId,
		"id":           c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.courseWork.studentSubmissions.return" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesCourseWorkStudentSubmissionsReturnCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a student submission.\n\nReturning a student submission transfers ownership of attached Drive\nfiles to the student and may also update the submission state.\nUnlike the Classroom application, returning a student submission does not\nset assignedGrade to the draftGrade value.\n\nOnly a teacher of the course that contains the requested student submission\nmay call this method.\n\nThis request must be made by the Developer Console project of the\n[OAuth client ID](https://support.google.com/cloud/answer/6158849) used to\ncreate the corresponding course work item.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access the\nrequested course or course work, return the requested student submission,\nor for access errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course, course work, or student submission\ndoes not exist.",
	//   "flatPath": "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}:return",
	//   "httpMethod": "POST",
	//   "id": "classroom.courses.courseWork.studentSubmissions.return",
	//   "parameterOrder": [
	//     "courseId",
	//     "courseWorkId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "courseWorkId": {
	//       "description": "Identifier of the course work.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "description": "Identifier of the student submission.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}:return",
	//   "request": {
	//     "$ref": "ReturnStudentSubmissionRequest"
	//   },
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.coursework.students"
	//   ]
	// }

}

// method id "classroom.courses.courseWork.studentSubmissions.turnIn":

type CoursesCourseWorkStudentSubmissionsTurnInCall struct {
	s                              *Service
	courseId                       string
	courseWorkId                   string
	id                             string
	turninstudentsubmissionrequest *TurnInStudentSubmissionRequest
	urlParams_                     gensupport.URLParams
	ctx_                           context.Context
	header_                        http.Header
}

// TurnIn: Turns in a student submission.
//
// Turning in a student submission transfers ownership of attached
// Drive
// files to the teacher and may also update the submission state.
//
// This may only be called by the student that owns the specified
// student
// submission.
//
// This request must be made by the Developer Console project of
// the
// [OAuth client ID](https://support.google.com/cloud/answer/6158849)
// used to
// create the corresponding course work item.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access the
// requested course or course work, turn in the requested student
// submission,
// or for access errors.
// * `INVALID_ARGUMENT` if the request is malformed.
// * `NOT_FOUND` if the requested course, course work, or student
// submission
// does not exist.
func (r *CoursesCourseWorkStudentSubmissionsService) TurnIn(courseId string, courseWorkId string, id string, turninstudentsubmissionrequest *TurnInStudentSubmissionRequest) *CoursesCourseWorkStudentSubmissionsTurnInCall {
	c := &CoursesCourseWorkStudentSubmissionsTurnInCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.courseWorkId = courseWorkId
	c.id = id
	c.turninstudentsubmissionrequest = turninstudentsubmissionrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesCourseWorkStudentSubmissionsTurnInCall) Fields(s ...googleapi.Field) *CoursesCourseWorkStudentSubmissionsTurnInCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesCourseWorkStudentSubmissionsTurnInCall) Context(ctx context.Context) *CoursesCourseWorkStudentSubmissionsTurnInCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesCourseWorkStudentSubmissionsTurnInCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesCourseWorkStudentSubmissionsTurnInCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.turninstudentsubmissionrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}:turnIn")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId":     c.courseId,
		"courseWorkId": c.courseWorkId,
		"id":           c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.courseWork.studentSubmissions.turnIn" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesCourseWorkStudentSubmissionsTurnInCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Turns in a student submission.\n\nTurning in a student submission transfers ownership of attached Drive\nfiles to the teacher and may also update the submission state.\n\nThis may only be called by the student that owns the specified student\nsubmission.\n\nThis request must be made by the Developer Console project of the\n[OAuth client ID](https://support.google.com/cloud/answer/6158849) used to\ncreate the corresponding course work item.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access the\nrequested course or course work, turn in the requested student submission,\nor for access errors.\n* `INVALID_ARGUMENT` if the request is malformed.\n* `NOT_FOUND` if the requested course, course work, or student submission\ndoes not exist.",
	//   "flatPath": "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}:turnIn",
	//   "httpMethod": "POST",
	//   "id": "classroom.courses.courseWork.studentSubmissions.turnIn",
	//   "parameterOrder": [
	//     "courseId",
	//     "courseWorkId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "courseWorkId": {
	//       "description": "Identifier of the course work.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "description": "Identifier of the student submission.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/courseWork/{courseWorkId}/studentSubmissions/{id}:turnIn",
	//   "request": {
	//     "$ref": "TurnInStudentSubmissionRequest"
	//   },
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.coursework.me"
	//   ]
	// }

}

// method id "classroom.courses.students.create":

type CoursesStudentsCreateCall struct {
	s          *Service
	courseId   string
	student    *Student
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Create: Adds a user as a student of a course.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// create
// students in this course or for access errors.
// * `NOT_FOUND` if the requested course ID does not exist.
// * `FAILED_PRECONDITION` if the requested user's account is
// disabled,
// for the following request errors:
//     * CourseMemberLimitReached
//     * CourseNotModifiable
//     * UserGroupsMembershipLimitReached
// * `ALREADY_EXISTS` if the user is already a student or teacher in
// the
// course.
func (r *CoursesStudentsService) Create(courseId string, student *Student) *CoursesStudentsCreateCall {
	c := &CoursesStudentsCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.student = student
	return c
}

// EnrollmentCode sets the optional parameter "enrollmentCode":
// Enrollment code of the course to create the student in.
// This code is required if userId
// corresponds to the requesting user; it may be omitted if the
// requesting
// user has administrative permissions to create students for any user.
func (c *CoursesStudentsCreateCall) EnrollmentCode(enrollmentCode string) *CoursesStudentsCreateCall {
	c.urlParams_.Set("enrollmentCode", enrollmentCode)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesStudentsCreateCall) Fields(s ...googleapi.Field) *CoursesStudentsCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesStudentsCreateCall) Context(ctx context.Context) *CoursesStudentsCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesStudentsCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesStudentsCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.student)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/students")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.students.create" call.
// Exactly one of *Student or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Student.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesStudentsCreateCall) Do(opts ...googleapi.CallOption) (*Student, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Student{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Adds a user as a student of a course.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to create\nstudents in this course or for access errors.\n* `NOT_FOUND` if the requested course ID does not exist.\n* `FAILED_PRECONDITION` if the requested user's account is disabled,\nfor the following request errors:\n    * CourseMemberLimitReached\n    * CourseNotModifiable\n    * UserGroupsMembershipLimitReached\n* `ALREADY_EXISTS` if the user is already a student or teacher in the\ncourse.",
	//   "flatPath": "v1/courses/{courseId}/students",
	//   "httpMethod": "POST",
	//   "id": "classroom.courses.students.create",
	//   "parameterOrder": [
	//     "courseId"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course to create the student in.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "enrollmentCode": {
	//       "description": "Enrollment code of the course to create the student in.\nThis code is required if userId\ncorresponds to the requesting user; it may be omitted if the requesting\nuser has administrative permissions to create students for any user.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/students",
	//   "request": {
	//     "$ref": "Student"
	//   },
	//   "response": {
	//     "$ref": "Student"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.profile.emails",
	//     "https://www.googleapis.com/auth/classroom.profile.photos",
	//     "https://www.googleapis.com/auth/classroom.rosters"
	//   ]
	// }

}

// method id "classroom.courses.students.delete":

type CoursesStudentsDeleteCall struct {
	s          *Service
	courseId   string
	userId     string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes a student of a course.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// delete
// students of this course or for access errors.
// * `NOT_FOUND` if no student of this course has the requested ID or if
// the
// course does not exist.
func (r *CoursesStudentsService) Delete(courseId string, userId string) *CoursesStudentsDeleteCall {
	c := &CoursesStudentsDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.userId = userId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesStudentsDeleteCall) Fields(s ...googleapi.Field) *CoursesStudentsDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesStudentsDeleteCall) Context(ctx context.Context) *CoursesStudentsDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesStudentsDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesStudentsDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/students/{userId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
		"userId":   c.userId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.students.delete" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesStudentsDeleteCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Deletes a student of a course.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to delete\nstudents of this course or for access errors.\n* `NOT_FOUND` if no student of this course has the requested ID or if the\ncourse does not exist.",
	//   "flatPath": "v1/courses/{courseId}/students/{userId}",
	//   "httpMethod": "DELETE",
	//   "id": "classroom.courses.students.delete",
	//   "parameterOrder": [
	//     "courseId",
	//     "userId"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "userId": {
	//       "description": "Identifier of the student to delete. The identifier can be one of the\nfollowing:\n\n* the numeric identifier for the user\n* the email address of the user\n* the string literal `\"me\"`, indicating the requesting user",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/students/{userId}",
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.rosters"
	//   ]
	// }

}

// method id "classroom.courses.students.get":

type CoursesStudentsGetCall struct {
	s            *Service
	courseId     string
	userId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Returns a student of a course.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// view
// students of this course or for access errors.
// * `NOT_FOUND` if no student of this course has the requested ID or if
// the
// course does not exist.
func (r *CoursesStudentsService) Get(courseId string, userId string) *CoursesStudentsGetCall {
	c := &CoursesStudentsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.userId = userId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesStudentsGetCall) Fields(s ...googleapi.Field) *CoursesStudentsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *CoursesStudentsGetCall) IfNoneMatch(entityTag string) *CoursesStudentsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesStudentsGetCall) Context(ctx context.Context) *CoursesStudentsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesStudentsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesStudentsGetCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/students/{userId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
		"userId":   c.userId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.students.get" call.
// Exactly one of *Student or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Student.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesStudentsGetCall) Do(opts ...googleapi.CallOption) (*Student, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Student{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a student of a course.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to view\nstudents of this course or for access errors.\n* `NOT_FOUND` if no student of this course has the requested ID or if the\ncourse does not exist.",
	//   "flatPath": "v1/courses/{courseId}/students/{userId}",
	//   "httpMethod": "GET",
	//   "id": "classroom.courses.students.get",
	//   "parameterOrder": [
	//     "courseId",
	//     "userId"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "userId": {
	//       "description": "Identifier of the student to return. The identifier can be one of the\nfollowing:\n\n* the numeric identifier for the user\n* the email address of the user\n* the string literal `\"me\"`, indicating the requesting user",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/students/{userId}",
	//   "response": {
	//     "$ref": "Student"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.profile.emails",
	//     "https://www.googleapis.com/auth/classroom.profile.photos",
	//     "https://www.googleapis.com/auth/classroom.rosters",
	//     "https://www.googleapis.com/auth/classroom.rosters.readonly"
	//   ]
	// }

}

// method id "classroom.courses.students.list":

type CoursesStudentsListCall struct {
	s            *Service
	courseId     string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Returns a list of students of this course that the requester
// is permitted to view.
//
// This method returns the following error codes:
//
// * `NOT_FOUND` if the course does not exist.
// * `PERMISSION_DENIED` for access errors.
func (r *CoursesStudentsService) List(courseId string) *CoursesStudentsListCall {
	c := &CoursesStudentsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum number of
// items to return. Zero means no maximum.
//
// The server may return fewer than the specified number of results.
func (c *CoursesStudentsListCall) PageSize(pageSize int64) *CoursesStudentsListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken":
// nextPageToken
// value returned from a previous
// list call, indicating that
// the subsequent page of results should be returned.
//
// The list request must be
// otherwise identical to the one that resulted in this token.
func (c *CoursesStudentsListCall) PageToken(pageToken string) *CoursesStudentsListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesStudentsListCall) Fields(s ...googleapi.Field) *CoursesStudentsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *CoursesStudentsListCall) IfNoneMatch(entityTag string) *CoursesStudentsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesStudentsListCall) Context(ctx context.Context) *CoursesStudentsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesStudentsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesStudentsListCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/students")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.students.list" call.
// Exactly one of *ListStudentsResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListStudentsResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *CoursesStudentsListCall) Do(opts ...googleapi.CallOption) (*ListStudentsResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListStudentsResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a list of students of this course that the requester\nis permitted to view.\n\nThis method returns the following error codes:\n\n* `NOT_FOUND` if the course does not exist.\n* `PERMISSION_DENIED` for access errors.",
	//   "flatPath": "v1/courses/{courseId}/students",
	//   "httpMethod": "GET",
	//   "id": "classroom.courses.students.list",
	//   "parameterOrder": [
	//     "courseId"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum number of items to return. Zero means no maximum.\n\nThe server may return fewer than the specified number of results.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "nextPageToken\nvalue returned from a previous\nlist call, indicating that\nthe subsequent page of results should be returned.\n\nThe list request must be\notherwise identical to the one that resulted in this token.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/students",
	//   "response": {
	//     "$ref": "ListStudentsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.profile.emails",
	//     "https://www.googleapis.com/auth/classroom.profile.photos",
	//     "https://www.googleapis.com/auth/classroom.rosters",
	//     "https://www.googleapis.com/auth/classroom.rosters.readonly"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *CoursesStudentsListCall) Pages(ctx context.Context, f func(*ListStudentsResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "classroom.courses.teachers.create":

type CoursesTeachersCreateCall struct {
	s          *Service
	courseId   string
	teacher    *Teacher
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Create: Creates a teacher of a course.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not  permitted to
// create
// teachers in this course or for access errors.
// * `NOT_FOUND` if the requested course ID does not exist.
// * `FAILED_PRECONDITION` if the requested user's account is
// disabled,
// for the following request errors:
//     * CourseMemberLimitReached
//     * CourseNotModifiable
//     * CourseTeacherLimitReached
//     * UserGroupsMembershipLimitReached
// * `ALREADY_EXISTS` if the user is already a teacher or student in
// the
// course.
func (r *CoursesTeachersService) Create(courseId string, teacher *Teacher) *CoursesTeachersCreateCall {
	c := &CoursesTeachersCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.teacher = teacher
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesTeachersCreateCall) Fields(s ...googleapi.Field) *CoursesTeachersCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesTeachersCreateCall) Context(ctx context.Context) *CoursesTeachersCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesTeachersCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesTeachersCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.teacher)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/teachers")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.teachers.create" call.
// Exactly one of *Teacher or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Teacher.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesTeachersCreateCall) Do(opts ...googleapi.CallOption) (*Teacher, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Teacher{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Creates a teacher of a course.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not  permitted to create\nteachers in this course or for access errors.\n* `NOT_FOUND` if the requested course ID does not exist.\n* `FAILED_PRECONDITION` if the requested user's account is disabled,\nfor the following request errors:\n    * CourseMemberLimitReached\n    * CourseNotModifiable\n    * CourseTeacherLimitReached\n    * UserGroupsMembershipLimitReached\n* `ALREADY_EXISTS` if the user is already a teacher or student in the\ncourse.",
	//   "flatPath": "v1/courses/{courseId}/teachers",
	//   "httpMethod": "POST",
	//   "id": "classroom.courses.teachers.create",
	//   "parameterOrder": [
	//     "courseId"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/teachers",
	//   "request": {
	//     "$ref": "Teacher"
	//   },
	//   "response": {
	//     "$ref": "Teacher"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.profile.emails",
	//     "https://www.googleapis.com/auth/classroom.profile.photos",
	//     "https://www.googleapis.com/auth/classroom.rosters"
	//   ]
	// }

}

// method id "classroom.courses.teachers.delete":

type CoursesTeachersDeleteCall struct {
	s          *Service
	courseId   string
	userId     string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes a teacher of a course.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// delete
// teachers of this course or for access errors.
// * `NOT_FOUND` if no teacher of this course has the requested ID or if
// the
// course does not exist.
// * `FAILED_PRECONDITION` if the requested ID belongs to the primary
// teacher
// of this course.
func (r *CoursesTeachersService) Delete(courseId string, userId string) *CoursesTeachersDeleteCall {
	c := &CoursesTeachersDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.userId = userId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesTeachersDeleteCall) Fields(s ...googleapi.Field) *CoursesTeachersDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesTeachersDeleteCall) Context(ctx context.Context) *CoursesTeachersDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesTeachersDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesTeachersDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/teachers/{userId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
		"userId":   c.userId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.teachers.delete" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesTeachersDeleteCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Deletes a teacher of a course.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to delete\nteachers of this course or for access errors.\n* `NOT_FOUND` if no teacher of this course has the requested ID or if the\ncourse does not exist.\n* `FAILED_PRECONDITION` if the requested ID belongs to the primary teacher\nof this course.",
	//   "flatPath": "v1/courses/{courseId}/teachers/{userId}",
	//   "httpMethod": "DELETE",
	//   "id": "classroom.courses.teachers.delete",
	//   "parameterOrder": [
	//     "courseId",
	//     "userId"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "userId": {
	//       "description": "Identifier of the teacher to delete. The identifier can be one of the\nfollowing:\n\n* the numeric identifier for the user\n* the email address of the user\n* the string literal `\"me\"`, indicating the requesting user",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/teachers/{userId}",
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.rosters"
	//   ]
	// }

}

// method id "classroom.courses.teachers.get":

type CoursesTeachersGetCall struct {
	s            *Service
	courseId     string
	userId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Returns a teacher of a course.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// view
// teachers of this course or for access errors.
// * `NOT_FOUND` if no teacher of this course has the requested ID or if
// the
// course does not exist.
func (r *CoursesTeachersService) Get(courseId string, userId string) *CoursesTeachersGetCall {
	c := &CoursesTeachersGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	c.userId = userId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesTeachersGetCall) Fields(s ...googleapi.Field) *CoursesTeachersGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *CoursesTeachersGetCall) IfNoneMatch(entityTag string) *CoursesTeachersGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesTeachersGetCall) Context(ctx context.Context) *CoursesTeachersGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesTeachersGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesTeachersGetCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/teachers/{userId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
		"userId":   c.userId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.teachers.get" call.
// Exactly one of *Teacher or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Teacher.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *CoursesTeachersGetCall) Do(opts ...googleapi.CallOption) (*Teacher, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Teacher{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a teacher of a course.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to view\nteachers of this course or for access errors.\n* `NOT_FOUND` if no teacher of this course has the requested ID or if the\ncourse does not exist.",
	//   "flatPath": "v1/courses/{courseId}/teachers/{userId}",
	//   "httpMethod": "GET",
	//   "id": "classroom.courses.teachers.get",
	//   "parameterOrder": [
	//     "courseId",
	//     "userId"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "userId": {
	//       "description": "Identifier of the teacher to return. The identifier can be one of the\nfollowing:\n\n* the numeric identifier for the user\n* the email address of the user\n* the string literal `\"me\"`, indicating the requesting user",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/teachers/{userId}",
	//   "response": {
	//     "$ref": "Teacher"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.profile.emails",
	//     "https://www.googleapis.com/auth/classroom.profile.photos",
	//     "https://www.googleapis.com/auth/classroom.rosters",
	//     "https://www.googleapis.com/auth/classroom.rosters.readonly"
	//   ]
	// }

}

// method id "classroom.courses.teachers.list":

type CoursesTeachersListCall struct {
	s            *Service
	courseId     string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Returns a list of teachers of this course that the requester
// is permitted to view.
//
// This method returns the following error codes:
//
// * `NOT_FOUND` if the course does not exist.
// * `PERMISSION_DENIED` for access errors.
func (r *CoursesTeachersService) List(courseId string) *CoursesTeachersListCall {
	c := &CoursesTeachersListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.courseId = courseId
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum number of
// items to return. Zero means no maximum.
//
// The server may return fewer than the specified number of results.
func (c *CoursesTeachersListCall) PageSize(pageSize int64) *CoursesTeachersListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken":
// nextPageToken
// value returned from a previous
// list call, indicating that
// the subsequent page of results should be returned.
//
// The list request must be
// otherwise identical to the one that resulted in this token.
func (c *CoursesTeachersListCall) PageToken(pageToken string) *CoursesTeachersListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *CoursesTeachersListCall) Fields(s ...googleapi.Field) *CoursesTeachersListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *CoursesTeachersListCall) IfNoneMatch(entityTag string) *CoursesTeachersListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *CoursesTeachersListCall) Context(ctx context.Context) *CoursesTeachersListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *CoursesTeachersListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *CoursesTeachersListCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/courses/{courseId}/teachers")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"courseId": c.courseId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.courses.teachers.list" call.
// Exactly one of *ListTeachersResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListTeachersResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *CoursesTeachersListCall) Do(opts ...googleapi.CallOption) (*ListTeachersResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListTeachersResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a list of teachers of this course that the requester\nis permitted to view.\n\nThis method returns the following error codes:\n\n* `NOT_FOUND` if the course does not exist.\n* `PERMISSION_DENIED` for access errors.",
	//   "flatPath": "v1/courses/{courseId}/teachers",
	//   "httpMethod": "GET",
	//   "id": "classroom.courses.teachers.list",
	//   "parameterOrder": [
	//     "courseId"
	//   ],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Identifier of the course.\nThis identifier can be either the Classroom-assigned identifier or an\nalias.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum number of items to return. Zero means no maximum.\n\nThe server may return fewer than the specified number of results.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "nextPageToken\nvalue returned from a previous\nlist call, indicating that\nthe subsequent page of results should be returned.\n\nThe list request must be\notherwise identical to the one that resulted in this token.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/courses/{courseId}/teachers",
	//   "response": {
	//     "$ref": "ListTeachersResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.profile.emails",
	//     "https://www.googleapis.com/auth/classroom.profile.photos",
	//     "https://www.googleapis.com/auth/classroom.rosters",
	//     "https://www.googleapis.com/auth/classroom.rosters.readonly"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *CoursesTeachersListCall) Pages(ctx context.Context, f func(*ListTeachersResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "classroom.invitations.accept":

type InvitationsAcceptCall struct {
	s          *Service
	id         string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Accept: Accepts an invitation, removing it and adding the invited
// user to the
// teachers or students (as appropriate) of the specified course. Only
// the
// invited user may accept an invitation.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// accept the
// requested invitation or for access errors.
// * `FAILED_PRECONDITION` for the following request errors:
//     * CourseMemberLimitReached
//     * CourseNotModifiable
//     * CourseTeacherLimitReached
//     * UserGroupsMembershipLimitReached
// * `NOT_FOUND` if no invitation exists with the requested ID.
func (r *InvitationsService) Accept(id string) *InvitationsAcceptCall {
	c := &InvitationsAcceptCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.id = id
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *InvitationsAcceptCall) Fields(s ...googleapi.Field) *InvitationsAcceptCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *InvitationsAcceptCall) Context(ctx context.Context) *InvitationsAcceptCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *InvitationsAcceptCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *InvitationsAcceptCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/invitations/{id}:accept")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"id": c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.invitations.accept" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *InvitationsAcceptCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Accepts an invitation, removing it and adding the invited user to the\nteachers or students (as appropriate) of the specified course. Only the\ninvited user may accept an invitation.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to accept the\nrequested invitation or for access errors.\n* `FAILED_PRECONDITION` for the following request errors:\n    * CourseMemberLimitReached\n    * CourseNotModifiable\n    * CourseTeacherLimitReached\n    * UserGroupsMembershipLimitReached\n* `NOT_FOUND` if no invitation exists with the requested ID.",
	//   "flatPath": "v1/invitations/{id}:accept",
	//   "httpMethod": "POST",
	//   "id": "classroom.invitations.accept",
	//   "parameterOrder": [
	//     "id"
	//   ],
	//   "parameters": {
	//     "id": {
	//       "description": "Identifier of the invitation to accept.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/invitations/{id}:accept",
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.rosters"
	//   ]
	// }

}

// method id "classroom.invitations.create":

type InvitationsCreateCall struct {
	s          *Service
	invitation *Invitation
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Create: Creates an invitation. Only one invitation for a user and
// course may exist
// at a time. Delete and re-create an invitation to make changes.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// create
// invitations for this course or for access errors.
// * `NOT_FOUND` if the course or the user does not exist.
// * `FAILED_PRECONDITION` if the requested user's account is disabled
// or if
// the user already has this role or a role with greater permissions.
// * `ALREADY_EXISTS` if an invitation for the specified user and
// course
// already exists.
func (r *InvitationsService) Create(invitation *Invitation) *InvitationsCreateCall {
	c := &InvitationsCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.invitation = invitation
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *InvitationsCreateCall) Fields(s ...googleapi.Field) *InvitationsCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *InvitationsCreateCall) Context(ctx context.Context) *InvitationsCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *InvitationsCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *InvitationsCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.invitation)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/invitations")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.invitations.create" call.
// Exactly one of *Invitation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Invitation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *InvitationsCreateCall) Do(opts ...googleapi.CallOption) (*Invitation, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Invitation{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Creates an invitation. Only one invitation for a user and course may exist\nat a time. Delete and re-create an invitation to make changes.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to create\ninvitations for this course or for access errors.\n* `NOT_FOUND` if the course or the user does not exist.\n* `FAILED_PRECONDITION` if the requested user's account is disabled or if\nthe user already has this role or a role with greater permissions.\n* `ALREADY_EXISTS` if an invitation for the specified user and course\nalready exists.",
	//   "flatPath": "v1/invitations",
	//   "httpMethod": "POST",
	//   "id": "classroom.invitations.create",
	//   "parameterOrder": [],
	//   "parameters": {},
	//   "path": "v1/invitations",
	//   "request": {
	//     "$ref": "Invitation"
	//   },
	//   "response": {
	//     "$ref": "Invitation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.rosters"
	//   ]
	// }

}

// method id "classroom.invitations.delete":

type InvitationsDeleteCall struct {
	s          *Service
	id         string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes an invitation.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// delete the
// requested invitation or for access errors.
// * `NOT_FOUND` if no invitation exists with the requested ID.
func (r *InvitationsService) Delete(id string) *InvitationsDeleteCall {
	c := &InvitationsDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.id = id
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *InvitationsDeleteCall) Fields(s ...googleapi.Field) *InvitationsDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *InvitationsDeleteCall) Context(ctx context.Context) *InvitationsDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *InvitationsDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *InvitationsDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/invitations/{id}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"id": c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.invitations.delete" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *InvitationsDeleteCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Deletes an invitation.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to delete the\nrequested invitation or for access errors.\n* `NOT_FOUND` if no invitation exists with the requested ID.",
	//   "flatPath": "v1/invitations/{id}",
	//   "httpMethod": "DELETE",
	//   "id": "classroom.invitations.delete",
	//   "parameterOrder": [
	//     "id"
	//   ],
	//   "parameters": {
	//     "id": {
	//       "description": "Identifier of the invitation to delete.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/invitations/{id}",
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.rosters"
	//   ]
	// }

}

// method id "classroom.invitations.get":

type InvitationsGetCall struct {
	s            *Service
	id           string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Returns an invitation.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to view
// the
// requested invitation or for access errors.
// * `NOT_FOUND` if no invitation exists with the requested ID.
func (r *InvitationsService) Get(id string) *InvitationsGetCall {
	c := &InvitationsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.id = id
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *InvitationsGetCall) Fields(s ...googleapi.Field) *InvitationsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *InvitationsGetCall) IfNoneMatch(entityTag string) *InvitationsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *InvitationsGetCall) Context(ctx context.Context) *InvitationsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *InvitationsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *InvitationsGetCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/invitations/{id}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"id": c.id,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.invitations.get" call.
// Exactly one of *Invitation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Invitation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *InvitationsGetCall) Do(opts ...googleapi.CallOption) (*Invitation, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Invitation{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns an invitation.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to view the\nrequested invitation or for access errors.\n* `NOT_FOUND` if no invitation exists with the requested ID.",
	//   "flatPath": "v1/invitations/{id}",
	//   "httpMethod": "GET",
	//   "id": "classroom.invitations.get",
	//   "parameterOrder": [
	//     "id"
	//   ],
	//   "parameters": {
	//     "id": {
	//       "description": "Identifier of the invitation to return.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/invitations/{id}",
	//   "response": {
	//     "$ref": "Invitation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.rosters",
	//     "https://www.googleapis.com/auth/classroom.rosters.readonly"
	//   ]
	// }

}

// method id "classroom.invitations.list":

type InvitationsListCall struct {
	s            *Service
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Returns a list of invitations that the requesting user is
// permitted to
// view, restricted to those that match the list request.
//
// *Note:* At least one of `user_id` or `course_id` must be supplied.
// Both
// fields can be supplied.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` for access errors.
func (r *InvitationsService) List() *InvitationsListCall {
	c := &InvitationsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	return c
}

// CourseId sets the optional parameter "courseId": Restricts returned
// invitations to those for a course with the specified
// identifier.
func (c *InvitationsListCall) CourseId(courseId string) *InvitationsListCall {
	c.urlParams_.Set("courseId", courseId)
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum number of
// items to return. Zero means no maximum.
//
// The server may return fewer than the specified number of results.
func (c *InvitationsListCall) PageSize(pageSize int64) *InvitationsListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken":
// nextPageToken
// value returned from a previous
// list call, indicating
// that the subsequent page of results should be returned.
//
// The list request must be
// otherwise identical to the one that resulted in this token.
func (c *InvitationsListCall) PageToken(pageToken string) *InvitationsListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// UserId sets the optional parameter "userId": Restricts returned
// invitations to those for a specific user. The identifier
// can be one of the following:
//
// * the numeric identifier for the user
// * the email address of the user
// * the string literal "me", indicating the requesting user
func (c *InvitationsListCall) UserId(userId string) *InvitationsListCall {
	c.urlParams_.Set("userId", userId)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *InvitationsListCall) Fields(s ...googleapi.Field) *InvitationsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *InvitationsListCall) IfNoneMatch(entityTag string) *InvitationsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *InvitationsListCall) Context(ctx context.Context) *InvitationsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *InvitationsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *InvitationsListCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/invitations")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.invitations.list" call.
// Exactly one of *ListInvitationsResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListInvitationsResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *InvitationsListCall) Do(opts ...googleapi.CallOption) (*ListInvitationsResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListInvitationsResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a list of invitations that the requesting user is permitted to\nview, restricted to those that match the list request.\n\n*Note:* At least one of `user_id` or `course_id` must be supplied. Both\nfields can be supplied.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` for access errors.",
	//   "flatPath": "v1/invitations",
	//   "httpMethod": "GET",
	//   "id": "classroom.invitations.list",
	//   "parameterOrder": [],
	//   "parameters": {
	//     "courseId": {
	//       "description": "Restricts returned invitations to those for a course with the specified\nidentifier.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum number of items to return. Zero means no maximum.\n\nThe server may return fewer than the specified number of results.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "nextPageToken\nvalue returned from a previous\nlist call, indicating\nthat the subsequent page of results should be returned.\n\nThe list request must be\notherwise identical to the one that resulted in this token.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "userId": {
	//       "description": "Restricts returned invitations to those for a specific user. The identifier\ncan be one of the following:\n\n* the numeric identifier for the user\n* the email address of the user\n* the string literal `\"me\"`, indicating the requesting user",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/invitations",
	//   "response": {
	//     "$ref": "ListInvitationsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.rosters",
	//     "https://www.googleapis.com/auth/classroom.rosters.readonly"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *InvitationsListCall) Pages(ctx context.Context, f func(*ListInvitationsResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "classroom.registrations.create":

type RegistrationsCreateCall struct {
	s            *Service
	registration *Registration
	urlParams_   gensupport.URLParams
	ctx_         context.Context
	header_      http.Header
}

// Create: Creates a `Registration`, causing Classroom to start sending
// notifications
// from the provided `feed` to the provided `destination`.
//
// Returns the created `Registration`. Currently, this will be the same
// as
// the argument, but with server-assigned fields such as `expiry_time`
// and
// `id` filled in.
//
// Note that any value specified for the `expiry_time` or `id` fields
// will be
// ignored.
//
// While Classroom may validate the `destination` and return errors on a
// best
// effort basis, it is the caller's responsibility to ensure that it
// exists
// and that Classroom has permission to publish to it.
//
// This method may return the following error codes:
//
// * `PERMISSION_DENIED` if:
//     * the authenticated user does not have permission to receive
//       notifications from the requested field; or
//     * the credential provided does not include the appropriate scope
// for the
//       requested feed.
//     * another access error is encountered.
// * `INVALID_ARGUMENT` if:
//     * no `destination` is specified, or the specified `destination`
// is not
//       valid; or
//     * no `feed` is specified, or the specified `feed` is not valid.
// * `NOT_FOUND` if:
//     * the specified `feed` cannot be located, or the requesting user
// does not
//       have permission to determine whether or not it exists; or
//     * the specified `destination` cannot be located, or Classroom has
// not
//       been granted permission to publish to it.
func (r *RegistrationsService) Create(registration *Registration) *RegistrationsCreateCall {
	c := &RegistrationsCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.registration = registration
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *RegistrationsCreateCall) Fields(s ...googleapi.Field) *RegistrationsCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *RegistrationsCreateCall) Context(ctx context.Context) *RegistrationsCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *RegistrationsCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *RegistrationsCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.registration)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/registrations")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.registrations.create" call.
// Exactly one of *Registration or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Registration.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *RegistrationsCreateCall) Do(opts ...googleapi.CallOption) (*Registration, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Registration{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Creates a `Registration`, causing Classroom to start sending notifications\nfrom the provided `feed` to the provided `destination`.\n\nReturns the created `Registration`. Currently, this will be the same as\nthe argument, but with server-assigned fields such as `expiry_time` and\n`id` filled in.\n\nNote that any value specified for the `expiry_time` or `id` fields will be\nignored.\n\nWhile Classroom may validate the `destination` and return errors on a best\neffort basis, it is the caller's responsibility to ensure that it exists\nand that Classroom has permission to publish to it.\n\nThis method may return the following error codes:\n\n* `PERMISSION_DENIED` if:\n    * the authenticated user does not have permission to receive\n      notifications from the requested field; or\n    * the credential provided does not include the appropriate scope for the\n      requested feed.\n    * another access error is encountered.\n* `INVALID_ARGUMENT` if:\n    * no `destination` is specified, or the specified `destination` is not\n      valid; or\n    * no `feed` is specified, or the specified `feed` is not valid.\n* `NOT_FOUND` if:\n    * the specified `feed` cannot be located, or the requesting user does not\n      have permission to determine whether or not it exists; or\n    * the specified `destination` cannot be located, or Classroom has not\n      been granted permission to publish to it.",
	//   "flatPath": "v1/registrations",
	//   "httpMethod": "POST",
	//   "id": "classroom.registrations.create",
	//   "parameterOrder": [],
	//   "parameters": {},
	//   "path": "v1/registrations",
	//   "request": {
	//     "$ref": "Registration"
	//   },
	//   "response": {
	//     "$ref": "Registration"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.rosters",
	//     "https://www.googleapis.com/auth/classroom.rosters.readonly"
	//   ]
	// }

}

// method id "classroom.registrations.delete":

type RegistrationsDeleteCall struct {
	s              *Service
	registrationId string
	urlParams_     gensupport.URLParams
	ctx_           context.Context
	header_        http.Header
}

// Delete: Deletes a `Registration`, causing Classroom to stop sending
// notifications
// for that `Registration`.
func (r *RegistrationsService) Delete(registrationId string) *RegistrationsDeleteCall {
	c := &RegistrationsDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.registrationId = registrationId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *RegistrationsDeleteCall) Fields(s ...googleapi.Field) *RegistrationsDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *RegistrationsDeleteCall) Context(ctx context.Context) *RegistrationsDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *RegistrationsDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *RegistrationsDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/registrations/{registrationId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"registrationId": c.registrationId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.registrations.delete" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *RegistrationsDeleteCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Deletes a `Registration`, causing Classroom to stop sending notifications\nfor that `Registration`.",
	//   "flatPath": "v1/registrations/{registrationId}",
	//   "httpMethod": "DELETE",
	//   "id": "classroom.registrations.delete",
	//   "parameterOrder": [
	//     "registrationId"
	//   ],
	//   "parameters": {
	//     "registrationId": {
	//       "description": "The `registration_id` of the `Registration` to be deleted.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/registrations/{registrationId}",
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.rosters",
	//     "https://www.googleapis.com/auth/classroom.rosters.readonly"
	//   ]
	// }

}

// method id "classroom.userProfiles.get":

type UserProfilesGetCall struct {
	s            *Service
	userId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Returns a user profile.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// access
// this user profile, if no profile exists with the requested ID, or
// for
// access errors.
func (r *UserProfilesService) Get(userId string) *UserProfilesGetCall {
	c := &UserProfilesGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.userId = userId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UserProfilesGetCall) Fields(s ...googleapi.Field) *UserProfilesGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *UserProfilesGetCall) IfNoneMatch(entityTag string) *UserProfilesGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *UserProfilesGetCall) Context(ctx context.Context) *UserProfilesGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *UserProfilesGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *UserProfilesGetCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/userProfiles/{userId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"userId": c.userId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.userProfiles.get" call.
// Exactly one of *UserProfile or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *UserProfile.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *UserProfilesGetCall) Do(opts ...googleapi.CallOption) (*UserProfile, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &UserProfile{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a user profile.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to access\nthis user profile, if no profile exists with the requested ID, or for\naccess errors.",
	//   "flatPath": "v1/userProfiles/{userId}",
	//   "httpMethod": "GET",
	//   "id": "classroom.userProfiles.get",
	//   "parameterOrder": [
	//     "userId"
	//   ],
	//   "parameters": {
	//     "userId": {
	//       "description": "Identifier of the profile to return. The identifier can be one of the\nfollowing:\n\n* the numeric identifier for the user\n* the email address of the user\n* the string literal `\"me\"`, indicating the requesting user",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/userProfiles/{userId}",
	//   "response": {
	//     "$ref": "UserProfile"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.profile.emails",
	//     "https://www.googleapis.com/auth/classroom.profile.photos",
	//     "https://www.googleapis.com/auth/classroom.rosters",
	//     "https://www.googleapis.com/auth/classroom.rosters.readonly"
	//   ]
	// }

}

// method id "classroom.userProfiles.guardianInvitations.create":

type UserProfilesGuardianInvitationsCreateCall struct {
	s                  *Service
	studentId          string
	guardianinvitation *GuardianInvitation
	urlParams_         gensupport.URLParams
	ctx_               context.Context
	header_            http.Header
}

// Create: Creates a guardian invitation, and sends an email to the
// guardian asking
// them to confirm that they are the student's guardian.
//
// Once the guardian accepts the invitation, their `state` will change
// to
// `COMPLETED` and they will start receiving guardian notifications.
// A
// `Guardian` resource will also be created to represent the active
// guardian.
//
// The request object must have the `student_id`
// and
// `invited_email_address` fields set. Failing to set these fields,
// or
// setting any other fields in the request, will result in an
// error.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the current user does not have permission
// to
//   manage guardians, if the guardian in question has already rejected
//   too many requests for that student, if guardians are not enabled
// for the
//   domain in question, or for other access errors.
// * `RESOURCE_EXHAUSTED` if the student or guardian has exceeded the
// guardian
//   link limit.
// * `INVALID_ARGUMENT` if the guardian email address is not valid (for
//   example, if it is too long), or if the format of the student ID
// provided
//   cannot be recognized (it is not an email address, nor a `user_id`
// from
//   this API). This error will also be returned if read-only fields are
// set,
//   or if the `state` field is set to to a value other than
// `PENDING`.
// * `NOT_FOUND` if the student ID provided is a valid student ID, but
//   Classroom has no record of that student.
// * `ALREADY_EXISTS` if there is already a pending guardian invitation
// for
//   the student and `invited_email_address` provided, or if the
// provided
//   `invited_email_address` matches the Google account of an existing
//   `Guardian` for this user.
func (r *UserProfilesGuardianInvitationsService) Create(studentId string, guardianinvitation *GuardianInvitation) *UserProfilesGuardianInvitationsCreateCall {
	c := &UserProfilesGuardianInvitationsCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.studentId = studentId
	c.guardianinvitation = guardianinvitation
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UserProfilesGuardianInvitationsCreateCall) Fields(s ...googleapi.Field) *UserProfilesGuardianInvitationsCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *UserProfilesGuardianInvitationsCreateCall) Context(ctx context.Context) *UserProfilesGuardianInvitationsCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *UserProfilesGuardianInvitationsCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *UserProfilesGuardianInvitationsCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.guardianinvitation)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/userProfiles/{studentId}/guardianInvitations")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"studentId": c.studentId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.userProfiles.guardianInvitations.create" call.
// Exactly one of *GuardianInvitation or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *GuardianInvitation.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *UserProfilesGuardianInvitationsCreateCall) Do(opts ...googleapi.CallOption) (*GuardianInvitation, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &GuardianInvitation{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Creates a guardian invitation, and sends an email to the guardian asking\nthem to confirm that they are the student's guardian.\n\nOnce the guardian accepts the invitation, their `state` will change to\n`COMPLETED` and they will start receiving guardian notifications. A\n`Guardian` resource will also be created to represent the active guardian.\n\nThe request object must have the `student_id` and\n`invited_email_address` fields set. Failing to set these fields, or\nsetting any other fields in the request, will result in an error.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the current user does not have permission to\n  manage guardians, if the guardian in question has already rejected\n  too many requests for that student, if guardians are not enabled for the\n  domain in question, or for other access errors.\n* `RESOURCE_EXHAUSTED` if the student or guardian has exceeded the guardian\n  link limit.\n* `INVALID_ARGUMENT` if the guardian email address is not valid (for\n  example, if it is too long), or if the format of the student ID provided\n  cannot be recognized (it is not an email address, nor a `user_id` from\n  this API). This error will also be returned if read-only fields are set,\n  or if the `state` field is set to to a value other than `PENDING`.\n* `NOT_FOUND` if the student ID provided is a valid student ID, but\n  Classroom has no record of that student.\n* `ALREADY_EXISTS` if there is already a pending guardian invitation for\n  the student and `invited_email_address` provided, or if the provided\n  `invited_email_address` matches the Google account of an existing\n  `Guardian` for this user.",
	//   "flatPath": "v1/userProfiles/{studentId}/guardianInvitations",
	//   "httpMethod": "POST",
	//   "id": "classroom.userProfiles.guardianInvitations.create",
	//   "parameterOrder": [
	//     "studentId"
	//   ],
	//   "parameters": {
	//     "studentId": {
	//       "description": "ID of the student (in standard format)",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/userProfiles/{studentId}/guardianInvitations",
	//   "request": {
	//     "$ref": "GuardianInvitation"
	//   },
	//   "response": {
	//     "$ref": "GuardianInvitation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.guardianlinks.students"
	//   ]
	// }

}

// method id "classroom.userProfiles.guardianInvitations.get":

type UserProfilesGuardianInvitationsGetCall struct {
	s            *Service
	studentId    string
	invitationId string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Returns a specific guardian invitation.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the requesting user is not permitted to
// view
//   guardian invitations for the student identified by the
// `student_id`, if
//   guardians are not enabled for the domain in question, or for other
//   access errors.
// * `INVALID_ARGUMENT` if a `student_id` is specified, but its format
// cannot
//   be recognized (it is not an email address, nor a `student_id` from
// the
//   API, nor the literal string `me`).
// * `NOT_FOUND` if Classroom cannot find any record of the given
// student or
//   `invitation_id`. May also be returned if the student exists, but
// the
//   requesting user does not have access to see that student.
func (r *UserProfilesGuardianInvitationsService) Get(studentId string, invitationId string) *UserProfilesGuardianInvitationsGetCall {
	c := &UserProfilesGuardianInvitationsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.studentId = studentId
	c.invitationId = invitationId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UserProfilesGuardianInvitationsGetCall) Fields(s ...googleapi.Field) *UserProfilesGuardianInvitationsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *UserProfilesGuardianInvitationsGetCall) IfNoneMatch(entityTag string) *UserProfilesGuardianInvitationsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *UserProfilesGuardianInvitationsGetCall) Context(ctx context.Context) *UserProfilesGuardianInvitationsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *UserProfilesGuardianInvitationsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *UserProfilesGuardianInvitationsGetCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/userProfiles/{studentId}/guardianInvitations/{invitationId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"studentId":    c.studentId,
		"invitationId": c.invitationId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.userProfiles.guardianInvitations.get" call.
// Exactly one of *GuardianInvitation or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *GuardianInvitation.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *UserProfilesGuardianInvitationsGetCall) Do(opts ...googleapi.CallOption) (*GuardianInvitation, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &GuardianInvitation{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a specific guardian invitation.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the requesting user is not permitted to view\n  guardian invitations for the student identified by the `student_id`, if\n  guardians are not enabled for the domain in question, or for other\n  access errors.\n* `INVALID_ARGUMENT` if a `student_id` is specified, but its format cannot\n  be recognized (it is not an email address, nor a `student_id` from the\n  API, nor the literal string `me`).\n* `NOT_FOUND` if Classroom cannot find any record of the given student or\n  `invitation_id`. May also be returned if the student exists, but the\n  requesting user does not have access to see that student.",
	//   "flatPath": "v1/userProfiles/{studentId}/guardianInvitations/{invitationId}",
	//   "httpMethod": "GET",
	//   "id": "classroom.userProfiles.guardianInvitations.get",
	//   "parameterOrder": [
	//     "studentId",
	//     "invitationId"
	//   ],
	//   "parameters": {
	//     "invitationId": {
	//       "description": "The `id` field of the `GuardianInvitation` being requested.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "studentId": {
	//       "description": "The ID of the student whose guardian invitation is being requested.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/userProfiles/{studentId}/guardianInvitations/{invitationId}",
	//   "response": {
	//     "$ref": "GuardianInvitation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.guardianlinks.students",
	//     "https://www.googleapis.com/auth/classroom.guardianlinks.students.readonly"
	//   ]
	// }

}

// method id "classroom.userProfiles.guardianInvitations.list":

type UserProfilesGuardianInvitationsListCall struct {
	s            *Service
	studentId    string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Returns a list of guardian invitations that the requesting user
// is
// permitted to view, filtered by the parameters provided.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if a `student_id` is specified, and the
// requesting
//   user is not permitted to view guardian invitations for that
// student, if
//   "-" is specified as the `student_id` and the user is not a
// domain
//   administrator, if guardians are not enabled for the domain in
// question,
//   or for other access errors.
// * `INVALID_ARGUMENT` if a `student_id` is specified, but its format
// cannot
//   be recognized (it is not an email address, nor a `student_id` from
// the
//   API, nor the literal string `me`). May also be returned if an
// invalid
//   `page_token` or `state` is provided.
// * `NOT_FOUND` if a `student_id` is specified, and its format can be
//   recognized, but Classroom has no record of that student.
func (r *UserProfilesGuardianInvitationsService) List(studentId string) *UserProfilesGuardianInvitationsListCall {
	c := &UserProfilesGuardianInvitationsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.studentId = studentId
	return c
}

// InvitedEmailAddress sets the optional parameter
// "invitedEmailAddress": If specified, only results with the specified
// `invited_email_address`
// will be returned.
func (c *UserProfilesGuardianInvitationsListCall) InvitedEmailAddress(invitedEmailAddress string) *UserProfilesGuardianInvitationsListCall {
	c.urlParams_.Set("invitedEmailAddress", invitedEmailAddress)
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum number of
// items to return. Zero or unspecified indicates that the
// server may assign a maximum.
//
// The server may return fewer than the specified number of results.
func (c *UserProfilesGuardianInvitationsListCall) PageSize(pageSize int64) *UserProfilesGuardianInvitationsListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken":
// nextPageToken
// value returned from a previous
// list call,
// indicating that the subsequent page of results should be
// returned.
//
// The list request
// must be otherwise identical to the one that resulted in this token.
func (c *UserProfilesGuardianInvitationsListCall) PageToken(pageToken string) *UserProfilesGuardianInvitationsListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// States sets the optional parameter "states": If specified, only
// results with the specified `state` values will be
// returned. Otherwise, results with a `state` of `PENDING` will be
// returned.
//
// Possible values:
//   "GUARDIAN_INVITATION_STATE_UNSPECIFIED"
//   "PENDING"
//   "COMPLETE"
func (c *UserProfilesGuardianInvitationsListCall) States(states ...string) *UserProfilesGuardianInvitationsListCall {
	c.urlParams_.SetMulti("states", append([]string{}, states...))
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UserProfilesGuardianInvitationsListCall) Fields(s ...googleapi.Field) *UserProfilesGuardianInvitationsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *UserProfilesGuardianInvitationsListCall) IfNoneMatch(entityTag string) *UserProfilesGuardianInvitationsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *UserProfilesGuardianInvitationsListCall) Context(ctx context.Context) *UserProfilesGuardianInvitationsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *UserProfilesGuardianInvitationsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *UserProfilesGuardianInvitationsListCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/userProfiles/{studentId}/guardianInvitations")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"studentId": c.studentId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.userProfiles.guardianInvitations.list" call.
// Exactly one of *ListGuardianInvitationsResponse or error will be
// non-nil. Any non-2xx status code is an error. Response headers are in
// either *ListGuardianInvitationsResponse.ServerResponse.Header or (if
// a response was returned at all) in error.(*googleapi.Error).Header.
// Use googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *UserProfilesGuardianInvitationsListCall) Do(opts ...googleapi.CallOption) (*ListGuardianInvitationsResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListGuardianInvitationsResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a list of guardian invitations that the requesting user is\npermitted to view, filtered by the parameters provided.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if a `student_id` is specified, and the requesting\n  user is not permitted to view guardian invitations for that student, if\n  `\"-\"` is specified as the `student_id` and the user is not a domain\n  administrator, if guardians are not enabled for the domain in question,\n  or for other access errors.\n* `INVALID_ARGUMENT` if a `student_id` is specified, but its format cannot\n  be recognized (it is not an email address, nor a `student_id` from the\n  API, nor the literal string `me`). May also be returned if an invalid\n  `page_token` or `state` is provided.\n* `NOT_FOUND` if a `student_id` is specified, and its format can be\n  recognized, but Classroom has no record of that student.",
	//   "flatPath": "v1/userProfiles/{studentId}/guardianInvitations",
	//   "httpMethod": "GET",
	//   "id": "classroom.userProfiles.guardianInvitations.list",
	//   "parameterOrder": [
	//     "studentId"
	//   ],
	//   "parameters": {
	//     "invitedEmailAddress": {
	//       "description": "If specified, only results with the specified `invited_email_address`\nwill be returned.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum number of items to return. Zero or unspecified indicates that the\nserver may assign a maximum.\n\nThe server may return fewer than the specified number of results.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "nextPageToken\nvalue returned from a previous\nlist call,\nindicating that the subsequent page of results should be returned.\n\nThe list request\nmust be otherwise identical to the one that resulted in this token.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "states": {
	//       "description": "If specified, only results with the specified `state` values will be\nreturned. Otherwise, results with a `state` of `PENDING` will be returned.",
	//       "enum": [
	//         "GUARDIAN_INVITATION_STATE_UNSPECIFIED",
	//         "PENDING",
	//         "COMPLETE"
	//       ],
	//       "location": "query",
	//       "repeated": true,
	//       "type": "string"
	//     },
	//     "studentId": {
	//       "description": "The ID of the student whose guardian invitations are to be returned.\nThe identifier can be one of the following:\n\n* the numeric identifier for the user\n* the email address of the user\n* the string literal `\"me\"`, indicating the requesting user\n* the string literal `\"-\"`, indicating that results should be returned for\n  all students that the requesting user is permitted to view guardian\n  invitations.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/userProfiles/{studentId}/guardianInvitations",
	//   "response": {
	//     "$ref": "ListGuardianInvitationsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.guardianlinks.students",
	//     "https://www.googleapis.com/auth/classroom.guardianlinks.students.readonly"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *UserProfilesGuardianInvitationsListCall) Pages(ctx context.Context, f func(*ListGuardianInvitationsResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "classroom.userProfiles.guardianInvitations.patch":

type UserProfilesGuardianInvitationsPatchCall struct {
	s                  *Service
	studentId          string
	invitationId       string
	guardianinvitation *GuardianInvitation
	urlParams_         gensupport.URLParams
	ctx_               context.Context
	header_            http.Header
}

// Patch: Modifies a guardian invitation.
//
// Currently, the only valid modification is to change the `state`
// from
// `PENDING` to `COMPLETE`. This has the effect of withdrawing the
// invitation.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if the current user does not have permission
// to
//   manage guardians, if guardians are not enabled for the domain in
// question
//   or for other access errors.
// * `FAILED_PRECONDITION` if the guardian link is not in the `PENDING`
// state.
// * `INVALID_ARGUMENT` if the format of the student ID provided
//   cannot be recognized (it is not an email address, nor a `user_id`
// from
//   this API), or if the passed `GuardianInvitation` has a `state`
// other than
//   `COMPLETE`, or if it modifies fields other than `state`.
// * `NOT_FOUND` if the student ID provided is a valid student ID, but
//   Classroom has no record of that student, or if the `id` field does
// not
//   refer to a guardian invitation known to Classroom.
func (r *UserProfilesGuardianInvitationsService) Patch(studentId string, invitationId string, guardianinvitation *GuardianInvitation) *UserProfilesGuardianInvitationsPatchCall {
	c := &UserProfilesGuardianInvitationsPatchCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.studentId = studentId
	c.invitationId = invitationId
	c.guardianinvitation = guardianinvitation
	return c
}

// UpdateMask sets the optional parameter "updateMask": Mask that
// identifies which fields on the course to update.
// This field is required to do an update. The update will fail if
// invalid
// fields are specified. The following fields are valid:
//
// * `state`
//
// When set in a query parameter, this field should be specified
// as
//
// `updateMask=<field1>,<field2>,...`
func (c *UserProfilesGuardianInvitationsPatchCall) UpdateMask(updateMask string) *UserProfilesGuardianInvitationsPatchCall {
	c.urlParams_.Set("updateMask", updateMask)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UserProfilesGuardianInvitationsPatchCall) Fields(s ...googleapi.Field) *UserProfilesGuardianInvitationsPatchCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *UserProfilesGuardianInvitationsPatchCall) Context(ctx context.Context) *UserProfilesGuardianInvitationsPatchCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *UserProfilesGuardianInvitationsPatchCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *UserProfilesGuardianInvitationsPatchCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.guardianinvitation)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/userProfiles/{studentId}/guardianInvitations/{invitationId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"studentId":    c.studentId,
		"invitationId": c.invitationId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.userProfiles.guardianInvitations.patch" call.
// Exactly one of *GuardianInvitation or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *GuardianInvitation.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *UserProfilesGuardianInvitationsPatchCall) Do(opts ...googleapi.CallOption) (*GuardianInvitation, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &GuardianInvitation{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Modifies a guardian invitation.\n\nCurrently, the only valid modification is to change the `state` from\n`PENDING` to `COMPLETE`. This has the effect of withdrawing the invitation.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if the current user does not have permission to\n  manage guardians, if guardians are not enabled for the domain in question\n  or for other access errors.\n* `FAILED_PRECONDITION` if the guardian link is not in the `PENDING` state.\n* `INVALID_ARGUMENT` if the format of the student ID provided\n  cannot be recognized (it is not an email address, nor a `user_id` from\n  this API), or if the passed `GuardianInvitation` has a `state` other than\n  `COMPLETE`, or if it modifies fields other than `state`.\n* `NOT_FOUND` if the student ID provided is a valid student ID, but\n  Classroom has no record of that student, or if the `id` field does not\n  refer to a guardian invitation known to Classroom.",
	//   "flatPath": "v1/userProfiles/{studentId}/guardianInvitations/{invitationId}",
	//   "httpMethod": "PATCH",
	//   "id": "classroom.userProfiles.guardianInvitations.patch",
	//   "parameterOrder": [
	//     "studentId",
	//     "invitationId"
	//   ],
	//   "parameters": {
	//     "invitationId": {
	//       "description": "The `id` field of the `GuardianInvitation` to be modified.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "studentId": {
	//       "description": "The ID of the student whose guardian invitation is to be modified.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "updateMask": {
	//       "description": "Mask that identifies which fields on the course to update.\nThis field is required to do an update. The update will fail if invalid\nfields are specified. The following fields are valid:\n\n* `state`\n\nWhen set in a query parameter, this field should be specified as\n\n`updateMask=\u003cfield1\u003e,\u003cfield2\u003e,...`",
	//       "format": "google-fieldmask",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/userProfiles/{studentId}/guardianInvitations/{invitationId}",
	//   "request": {
	//     "$ref": "GuardianInvitation"
	//   },
	//   "response": {
	//     "$ref": "GuardianInvitation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.guardianlinks.students"
	//   ]
	// }

}

// method id "classroom.userProfiles.guardians.delete":

type UserProfilesGuardiansDeleteCall struct {
	s          *Service
	studentId  string
	guardianId string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes a guardian.
//
// The guardian will no longer receive guardian notifications and the
// guardian
// will no longer be accessible via the API.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if no user that matches the provided
// `student_id`
//   is visible to the requesting user, if the requesting user is not
//   permitted to manage guardians for the student identified by the
//   `student_id`, if guardians are not enabled for the domain in
// question,
//   or for other access errors.
// * `INVALID_ARGUMENT` if a `student_id` is specified, but its format
// cannot
//   be recognized (it is not an email address, nor a `student_id` from
// the
//   API).
// * `NOT_FOUND` if the requesting user is permitted to modify guardians
// for
//   the requested `student_id`, but no `Guardian` record exists for
// that
//   student with the provided `guardian_id`.
func (r *UserProfilesGuardiansService) Delete(studentId string, guardianId string) *UserProfilesGuardiansDeleteCall {
	c := &UserProfilesGuardiansDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.studentId = studentId
	c.guardianId = guardianId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UserProfilesGuardiansDeleteCall) Fields(s ...googleapi.Field) *UserProfilesGuardiansDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *UserProfilesGuardiansDeleteCall) Context(ctx context.Context) *UserProfilesGuardiansDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *UserProfilesGuardiansDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *UserProfilesGuardiansDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/userProfiles/{studentId}/guardians/{guardianId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"studentId":  c.studentId,
		"guardianId": c.guardianId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.userProfiles.guardians.delete" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *UserProfilesGuardiansDeleteCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Deletes a guardian.\n\nThe guardian will no longer receive guardian notifications and the guardian\nwill no longer be accessible via the API.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if no user that matches the provided `student_id`\n  is visible to the requesting user, if the requesting user is not\n  permitted to manage guardians for the student identified by the\n  `student_id`, if guardians are not enabled for the domain in question,\n  or for other access errors.\n* `INVALID_ARGUMENT` if a `student_id` is specified, but its format cannot\n  be recognized (it is not an email address, nor a `student_id` from the\n  API).\n* `NOT_FOUND` if the requesting user is permitted to modify guardians for\n  the requested `student_id`, but no `Guardian` record exists for that\n  student with the provided `guardian_id`.",
	//   "flatPath": "v1/userProfiles/{studentId}/guardians/{guardianId}",
	//   "httpMethod": "DELETE",
	//   "id": "classroom.userProfiles.guardians.delete",
	//   "parameterOrder": [
	//     "studentId",
	//     "guardianId"
	//   ],
	//   "parameters": {
	//     "guardianId": {
	//       "description": "The `id` field from a `Guardian`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "studentId": {
	//       "description": "The student whose guardian is to be deleted. One of the following:\n\n* the numeric identifier for the user\n* the email address of the user\n* the string literal `\"me\"`, indicating the requesting user",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/userProfiles/{studentId}/guardians/{guardianId}",
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.guardianlinks.students"
	//   ]
	// }

}

// method id "classroom.userProfiles.guardians.get":

type UserProfilesGuardiansGetCall struct {
	s            *Service
	studentId    string
	guardianId   string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Returns a specific guardian.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if no user that matches the provided
// `student_id`
//   is visible to the requesting user, if the requesting user is not
//   permitted to view guardian information for the student identified
// by the
//   `student_id`, if guardians are not enabled for the domain in
// question,
//   or for other access errors.
// * `INVALID_ARGUMENT` if a `student_id` is specified, but its format
// cannot
//   be recognized (it is not an email address, nor a `student_id` from
// the
//   API, nor the literal string `me`).
// * `NOT_FOUND` if the requesting user is permitted to view guardians
// for
//   the requested `student_id`, but no `Guardian` record exists for
// that
//   student that matches the provided `guardian_id`.
func (r *UserProfilesGuardiansService) Get(studentId string, guardianId string) *UserProfilesGuardiansGetCall {
	c := &UserProfilesGuardiansGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.studentId = studentId
	c.guardianId = guardianId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UserProfilesGuardiansGetCall) Fields(s ...googleapi.Field) *UserProfilesGuardiansGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *UserProfilesGuardiansGetCall) IfNoneMatch(entityTag string) *UserProfilesGuardiansGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *UserProfilesGuardiansGetCall) Context(ctx context.Context) *UserProfilesGuardiansGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *UserProfilesGuardiansGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *UserProfilesGuardiansGetCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/userProfiles/{studentId}/guardians/{guardianId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"studentId":  c.studentId,
		"guardianId": c.guardianId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.userProfiles.guardians.get" call.
// Exactly one of *Guardian or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Guardian.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *UserProfilesGuardiansGetCall) Do(opts ...googleapi.CallOption) (*Guardian, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Guardian{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a specific guardian.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if no user that matches the provided `student_id`\n  is visible to the requesting user, if the requesting user is not\n  permitted to view guardian information for the student identified by the\n  `student_id`, if guardians are not enabled for the domain in question,\n  or for other access errors.\n* `INVALID_ARGUMENT` if a `student_id` is specified, but its format cannot\n  be recognized (it is not an email address, nor a `student_id` from the\n  API, nor the literal string `me`).\n* `NOT_FOUND` if the requesting user is permitted to view guardians for\n  the requested `student_id`, but no `Guardian` record exists for that\n  student that matches the provided `guardian_id`.",
	//   "flatPath": "v1/userProfiles/{studentId}/guardians/{guardianId}",
	//   "httpMethod": "GET",
	//   "id": "classroom.userProfiles.guardians.get",
	//   "parameterOrder": [
	//     "studentId",
	//     "guardianId"
	//   ],
	//   "parameters": {
	//     "guardianId": {
	//       "description": "The `id` field from a `Guardian`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "studentId": {
	//       "description": "The student whose guardian is being requested. One of the following:\n\n* the numeric identifier for the user\n* the email address of the user\n* the string literal `\"me\"`, indicating the requesting user",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/userProfiles/{studentId}/guardians/{guardianId}",
	//   "response": {
	//     "$ref": "Guardian"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.guardianlinks.me.readonly",
	//     "https://www.googleapis.com/auth/classroom.guardianlinks.students",
	//     "https://www.googleapis.com/auth/classroom.guardianlinks.students.readonly"
	//   ]
	// }

}

// method id "classroom.userProfiles.guardians.list":

type UserProfilesGuardiansListCall struct {
	s            *Service
	studentId    string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Returns a list of guardians that the requesting user is
// permitted to
// view, restricted to those that match the request.
//
// To list guardians for any student that the requesting user may
// view
// guardians for, use the literal character `-` for the student
// ID.
//
// This method returns the following error codes:
//
// * `PERMISSION_DENIED` if a `student_id` is specified, and the
// requesting
//   user is not permitted to view guardian information for that
// student, if
//   "-" is specified as the `student_id` and the user is not a
// domain
//   administrator, if guardians are not enabled for the domain in
// question,
//   if the `invited_email_address` filter is set by a user who is not
// a
//   domain administrator, or for other access errors.
// * `INVALID_ARGUMENT` if a `student_id` is specified, but its format
// cannot
//   be recognized (it is not an email address, nor a `student_id` from
// the
//   API, nor the literal string `me`). May also be returned if an
// invalid
//   `page_token` is provided.
// * `NOT_FOUND` if a `student_id` is specified, and its format can be
//   recognized, but Classroom has no record of that student.
func (r *UserProfilesGuardiansService) List(studentId string) *UserProfilesGuardiansListCall {
	c := &UserProfilesGuardiansListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.studentId = studentId
	return c
}

// InvitedEmailAddress sets the optional parameter
// "invitedEmailAddress": Filter results by the email address that the
// original invitation was sent
// to, resulting in this guardian link.
// This filter can only be used by domain administrators.
func (c *UserProfilesGuardiansListCall) InvitedEmailAddress(invitedEmailAddress string) *UserProfilesGuardiansListCall {
	c.urlParams_.Set("invitedEmailAddress", invitedEmailAddress)
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum number of
// items to return. Zero or unspecified indicates that the
// server may assign a maximum.
//
// The server may return fewer than the specified number of results.
func (c *UserProfilesGuardiansListCall) PageSize(pageSize int64) *UserProfilesGuardiansListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken":
// nextPageToken
// value returned from a previous
// list call,
// indicating that the subsequent page of results should be
// returned.
//
// The list request
// must be otherwise identical to the one that resulted in this token.
func (c *UserProfilesGuardiansListCall) PageToken(pageToken string) *UserProfilesGuardiansListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UserProfilesGuardiansListCall) Fields(s ...googleapi.Field) *UserProfilesGuardiansListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *UserProfilesGuardiansListCall) IfNoneMatch(entityTag string) *UserProfilesGuardiansListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *UserProfilesGuardiansListCall) Context(ctx context.Context) *UserProfilesGuardiansListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *UserProfilesGuardiansListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *UserProfilesGuardiansListCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/userProfiles/{studentId}/guardians")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"studentId": c.studentId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "classroom.userProfiles.guardians.list" call.
// Exactly one of *ListGuardiansResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListGuardiansResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *UserProfilesGuardiansListCall) Do(opts ...googleapi.CallOption) (*ListGuardiansResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListGuardiansResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Returns a list of guardians that the requesting user is permitted to\nview, restricted to those that match the request.\n\nTo list guardians for any student that the requesting user may view\nguardians for, use the literal character `-` for the student ID.\n\nThis method returns the following error codes:\n\n* `PERMISSION_DENIED` if a `student_id` is specified, and the requesting\n  user is not permitted to view guardian information for that student, if\n  `\"-\"` is specified as the `student_id` and the user is not a domain\n  administrator, if guardians are not enabled for the domain in question,\n  if the `invited_email_address` filter is set by a user who is not a\n  domain administrator, or for other access errors.\n* `INVALID_ARGUMENT` if a `student_id` is specified, but its format cannot\n  be recognized (it is not an email address, nor a `student_id` from the\n  API, nor the literal string `me`). May also be returned if an invalid\n  `page_token` is provided.\n* `NOT_FOUND` if a `student_id` is specified, and its format can be\n  recognized, but Classroom has no record of that student.",
	//   "flatPath": "v1/userProfiles/{studentId}/guardians",
	//   "httpMethod": "GET",
	//   "id": "classroom.userProfiles.guardians.list",
	//   "parameterOrder": [
	//     "studentId"
	//   ],
	//   "parameters": {
	//     "invitedEmailAddress": {
	//       "description": "Filter results by the email address that the original invitation was sent\nto, resulting in this guardian link.\nThis filter can only be used by domain administrators.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum number of items to return. Zero or unspecified indicates that the\nserver may assign a maximum.\n\nThe server may return fewer than the specified number of results.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "nextPageToken\nvalue returned from a previous\nlist call,\nindicating that the subsequent page of results should be returned.\n\nThe list request\nmust be otherwise identical to the one that resulted in this token.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "studentId": {
	//       "description": "Filter results by the student who the guardian is linked to.\nThe identifier can be one of the following:\n\n* the numeric identifier for the user\n* the email address of the user\n* the string literal `\"me\"`, indicating the requesting user\n* the string literal `\"-\"`, indicating that results should be returned for\n  all students that the requesting user has access to view.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/userProfiles/{studentId}/guardians",
	//   "response": {
	//     "$ref": "ListGuardiansResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/classroom.guardianlinks.me.readonly",
	//     "https://www.googleapis.com/auth/classroom.guardianlinks.students",
	//     "https://www.googleapis.com/auth/classroom.guardianlinks.students.readonly"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *UserProfilesGuardiansListCall) Pages(ctx context.Context, f func(*ListGuardiansResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}
