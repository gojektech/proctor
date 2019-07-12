package audit

import (
	"proctor/internal/app/proctord/storage"
	"proctor/internal/app/proctord/storage/postgres"
	"proctor/internal/app/service/infra/kubernetes"
	"testing"

	"proctor/internal/pkg/constant"
)

func TestJobsExecutionAuditing(t *testing.T) {
	mockStore := &storage.MockStore{}
	mockKubeClient := &kubernetes.MockKubernetesClient{}
	testAuditor := New(mockStore, mockKubeClient)
	jobsExecutionAuditLog := &postgres.JobsExecutionAuditLog{
		JobName: "any-job-name",
	}

	mockStore.On("AuditJobsExecution", jobsExecutionAuditLog).Return(nil).Once()

	testAuditor.JobsExecution(jobsExecutionAuditLog)

	mockStore.AssertExpectations(t)
	mockKubeClient.AssertExpectations(t)
}

func TestAuditJobsExecutionStatusAuditing(t *testing.T) {
	mockStore := &storage.MockStore{}
	mockKubeClient := &kubernetes.MockKubernetesClient{}
	testAuditor := New(mockStore, mockKubeClient)

	jobExecutionID := "job-execution-id"
	jobExecutionStatus := "job-execution-status"

	mockKubeClient.On("JobExecutionStatus", jobExecutionID).Return(jobExecutionStatus, nil)
	mockStore.On("UpdateJobsExecutionAuditLog", jobExecutionID, jobExecutionStatus).Return(nil).Once()

	_,_ = testAuditor.JobsExecutionStatus(jobExecutionID)

	mockStore.AssertExpectations(t)
	mockKubeClient.AssertExpectations(t)
}

func TestAuditJobsExecutionAndStatusAuditing(t *testing.T) {
	mockStore := &storage.MockStore{}
	mockKubeClient := &kubernetes.MockKubernetesClient{}
	testAuditor := New(mockStore, mockKubeClient)

	jobExecutionID := "job-execution-id"
	jobExecutionStatus := "job-execution-status"
	jobsExecutionAuditLog := &postgres.JobsExecutionAuditLog{
		JobName:             "any-job-name",
		ExecutionID:         postgres.StringToSQLString(jobExecutionID),
		JobSubmissionStatus: constant.JobSubmissionSuccess,
	}

	mockStore.On("AuditJobsExecution", jobsExecutionAuditLog).Return(nil).Once()

	mockKubeClient.On("JobExecutionStatus", jobExecutionID).Return(jobExecutionStatus, nil)
	mockStore.On("UpdateJobsExecutionAuditLog", jobExecutionID, jobExecutionStatus).Return(nil).Once()

	testAuditor.JobsExecutionAndStatus(jobsExecutionAuditLog)

	mockStore.AssertExpectations(t)
	mockKubeClient.AssertExpectations(t)
}
