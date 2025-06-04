/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package addons

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
)

func TestEKSAddonPlan(t *testing.T) {
	clusterName := "default.cluster"
	addonARN := "aws://someaddonarn"
	addon1Name := "addon1"
	addon1version := "1.0.0"
	addon1Upgrade := "2.0.0"
	addonStatusActive := string(ekstypes.AddonStatusActive)
	addonStatusUpdating := string(ekstypes.AddonStatusUpdating)
	addonStatusDeleting := string(ekstypes.AddonStatusDeleting)
	addonStatusCreating := string(ekstypes.AddonStatusCreating)
	created := time.Now()
	maxActiveUpdateDeleteWait := 30 * time.Minute

	testCases := []struct {
		name              string
		desiredAddons     []*EKSAddon
		installedAddons   []*EKSAddon
		expect            func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectCreateError bool
		expectDoError     bool
		preserveOnDelete  bool
	}{
		{
			name: "no desired and no installed",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				// Do nothing
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name: "no installed and 1 desired",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					CreateAddon(gomock.Eq(context.TODO()), gomock.Eq(&eks.CreateAddonInput{
						AddonName:        aws.String(addon1Name),
						AddonVersion:     aws.String(addon1version),
						ClusterName:      aws.String(clusterName),
						ResolveConflicts: ekstypes.ResolveConflictsOverwrite,
						Tags:             createTags(),
					})).
					Return(&eks.CreateAddonOutput{
						Addon: &ekstypes.Addon{
							AddonArn:     aws.String(addonARN),
							AddonName:    aws.String(addon1Name),
							AddonVersion: aws.String(addon1version),
							ClusterName:  aws.String(clusterName),
							CreatedAt:    &created,
							ModifiedAt:   &created,
							Status:       ekstypes.AddonStatusCreating,
							Tags:         createTags(),
						},
					}, nil)

				out := &eks.DescribeAddonOutput{
					Addon: &ekstypes.Addon{
						Status: ekstypes.AddonStatusActive,
					},
				}
				m.DescribeAddon(gomock.Eq(context.TODO()), gomock.Eq(&eks.DescribeAddonInput{
					AddonName:   aws.String(addon1Name),
					ClusterName: aws.String(clusterName),
				})).Return(out, nil)
			},
			desiredAddons: []*EKSAddon{
				createDesiredAddon(addon1Name, addon1version),
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name: "no installed and 2 desired",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					CreateAddon(gomock.Eq(context.TODO()), gomock.Eq(&eks.CreateAddonInput{
						AddonName:        aws.String(addon1Name),
						AddonVersion:     aws.String(addon1version),
						ClusterName:      aws.String(clusterName),
						ResolveConflicts: ekstypes.ResolveConflictsOverwrite,
						Tags:             createTags(),
					})).
					Return(&eks.CreateAddonOutput{
						Addon: &ekstypes.Addon{
							AddonArn:     aws.String(addonARN),
							AddonName:    aws.String(addon1Name),
							AddonVersion: aws.String(addon1version),
							ClusterName:  aws.String(clusterName),
							CreatedAt:    &created,
							ModifiedAt:   &created,
							Status:       ekstypes.AddonStatusCreating,
							Tags:         createTags(),
						},
					}, nil)

				out := &eks.DescribeAddonOutput{
					Addon: &ekstypes.Addon{
						Status: ekstypes.AddonStatusActive,
					},
				}
				m.DescribeAddon(gomock.Eq(context.TODO()), gomock.Eq(&eks.DescribeAddonInput{
					AddonName:   aws.String(addon1Name),
					ClusterName: aws.String(clusterName),
				})).Return(out, nil)

				m.
					CreateAddon(gomock.Eq(context.TODO()), gomock.Eq(&eks.CreateAddonInput{
						AddonName:        aws.String("addon2"),
						AddonVersion:     aws.String(addon1version),
						ClusterName:      aws.String(clusterName),
						ResolveConflicts: ekstypes.ResolveConflictsOverwrite,
						Tags:             createTags(),
					})).
					Return(&eks.CreateAddonOutput{
						Addon: &ekstypes.Addon{
							AddonArn:     aws.String(addonARN),
							AddonName:    aws.String("addon2"),
							AddonVersion: aws.String(addon1version),
							ClusterName:  aws.String(clusterName),
							CreatedAt:    &created,
							ModifiedAt:   &created,
							Status:       ekstypes.AddonStatusCreating,
							Tags:         createTags(),
						},
					}, nil)

				m.DescribeAddon(gomock.Eq(context.TODO()), gomock.Eq(&eks.DescribeAddonInput{
					AddonName:   aws.String("addon2"),
					ClusterName: aws.String(clusterName),
				})).Return(out, nil)
			},
			desiredAddons: []*EKSAddon{
				createDesiredAddon(addon1Name, addon1version),
				createDesiredAddon("addon2", addon1version),
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name: "1 installed and 1 desired - both same and installed active",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				// No Action expected
			},
			desiredAddons: []*EKSAddon{
				createDesiredAddon(addon1Name, addon1version),
			},
			installedAddons: []*EKSAddon{
				createInstalledAddon(addon1Name, addon1version, addonARN, addonStatusActive),
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name: "1 installed and 1 desired - both same and installed is creating",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				out := &eks.DescribeAddonOutput{
					Addon: &ekstypes.Addon{
						Status: ekstypes.AddonStatusActive,
					},
				}
				m.DescribeAddon(gomock.Eq(context.TODO()), gomock.Eq(&eks.DescribeAddonInput{
					AddonName:   aws.String(addon1Name),
					ClusterName: aws.String(clusterName),
				})).Return(out, nil)
			},
			desiredAddons: []*EKSAddon{
				createDesiredAddon(addon1Name, addon1version),
			},
			installedAddons: []*EKSAddon{
				createInstalledAddon(addon1Name, addon1version, addonARN, addonStatusCreating),
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name: "1 installed and 1 desired - version upgrade",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					UpdateAddon(gomock.Eq(context.TODO()), gomock.Eq(&eks.UpdateAddonInput{
						AddonName:        aws.String(addon1Name),
						AddonVersion:     aws.String(addon1Upgrade),
						ClusterName:      aws.String(clusterName),
						ResolveConflicts: ekstypes.ResolveConflictsOverwrite,
					})).
					Return(&eks.UpdateAddonOutput{
						Update: &ekstypes.Update{
							CreatedAt: &created,
							Id:        aws.String("someid"),
							Status:    ekstypes.UpdateStatus(ekstypes.AddonStatusUpdating),
							Type:      ekstypes.UpdateTypeVersionUpdate,
						},
					}, nil)

				out := &eks.DescribeAddonOutput{
					Addon: &ekstypes.Addon{
						Status: ekstypes.AddonStatusActive,
					},
				}
				m.DescribeAddon(gomock.Eq(context.TODO()), gomock.Eq(&eks.DescribeAddonInput{
					AddonName:   aws.String(addon1Name),
					ClusterName: aws.String(clusterName),
				})).Return(out, nil)
			},
			desiredAddons: []*EKSAddon{
				createDesiredAddon(addon1Name, addon1Upgrade),
			},
			installedAddons: []*EKSAddon{
				createInstalledAddon(addon1Name, addon1version, addonARN, addonStatusActive),
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name: "1 installed and 1 desired - version upgrade in progress",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				out := &eks.DescribeAddonOutput{
					Addon: &ekstypes.Addon{
						Status: ekstypes.AddonStatusActive,
					},
				}
				m.DescribeAddon(gomock.Eq(context.TODO()), gomock.Eq(&eks.DescribeAddonInput{
					AddonName:   aws.String(addon1Name),
					ClusterName: aws.String(clusterName),
				})).Return(out, nil)
			},
			desiredAddons: []*EKSAddon{
				createDesiredAddon(addon1Name, addon1Upgrade),
			},
			installedAddons: []*EKSAddon{
				createInstalledAddon(addon1Name, addon1Upgrade, addonARN, addonStatusUpdating),
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name: "1 installed and 1 desired - tags upgrade",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					TagResource(gomock.Eq(context.TODO()), gomock.Eq(&eks.TagResourceInput{
						ResourceArn: &addonARN,
						Tags:        createTagsAdditional(),
					})).
					Return(&eks.TagResourceOutput{}, nil)
			},
			desiredAddons: []*EKSAddon{
				createDesiredAddonExtraTag(addon1Name, addon1version),
			},
			installedAddons: []*EKSAddon{
				createInstalledAddon(addon1Name, addon1version, addonARN, addonStatusActive),
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name: "1 installed and 1 desired - version & tags upgrade",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					TagResource(gomock.Eq(context.TODO()), gomock.Eq(&eks.TagResourceInput{
						ResourceArn: &addonARN,
						Tags:        createTagsAdditional(),
					})).
					Return(&eks.TagResourceOutput{}, nil)
				m.
					UpdateAddon(gomock.Eq(context.TODO()), gomock.Eq(&eks.UpdateAddonInput{
						AddonName:        aws.String(addon1Name),
						AddonVersion:     aws.String(addon1Upgrade),
						ClusterName:      aws.String(clusterName),
						ResolveConflicts: ekstypes.ResolveConflictsOverwrite,
					})).
					Return(&eks.UpdateAddonOutput{
						Update: &ekstypes.Update{
							CreatedAt: &created,
							Id:        aws.String("someid"),
							Status:    ekstypes.UpdateStatus(ekstypes.AddonStatusUpdating),
							Type:      ekstypes.UpdateTypeVersionUpdate,
						},
					}, nil)

				out := &eks.DescribeAddonOutput{
					Addon: &ekstypes.Addon{
						Status: ekstypes.AddonStatusActive,
					},
				}
				m.DescribeAddon(gomock.Eq(context.TODO()), gomock.Eq(&eks.DescribeAddonInput{
					AddonName:   aws.String(addon1Name),
					ClusterName: aws.String(clusterName),
				})).Return(out, nil)
			},
			desiredAddons: []*EKSAddon{
				createDesiredAddonExtraTag(addon1Name, addon1Upgrade),
			},
			installedAddons: []*EKSAddon{
				createInstalledAddon(addon1Name, addon1version, addonARN, addonStatusActive),
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name: "1 installed and 0 desired - delete addon",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					DeleteAddon(gomock.Eq(context.TODO()), gomock.Eq(&eks.DeleteAddonInput{
						AddonName:   &addon1Name,
						ClusterName: &clusterName,
						Preserve:    false,
					})).
					Return(&eks.DeleteAddonOutput{
						Addon: &ekstypes.Addon{
							AddonArn:     aws.String(addonARN),
							AddonName:    aws.String(addon1Name),
							AddonVersion: aws.String(addon1version),
							ClusterName:  aws.String(clusterName),
							CreatedAt:    &created,
							ModifiedAt:   &created,
							Status:       ekstypes.AddonStatusDeleting,
							Tags:         createTags(),
						},
					}, nil)
				m.WaitUntilAddonDeleted(gomock.Eq(context.TODO()), gomock.Eq(&eks.DescribeAddonInput{
					AddonName:   aws.String(addon1Name),
					ClusterName: aws.String(clusterName),
				}), maxActiveUpdateDeleteWait).Return(nil)
			},
			installedAddons: []*EKSAddon{
				createInstalledAddon(addon1Name, addon1version, addonARN, addonStatusActive),
			},
			expectCreateError: false,
			expectDoError:     false,
			preserveOnDelete:  true,
		},
		{
			name: "1 installed and 0 desired - delete addon & preserve",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					DeleteAddon(gomock.Eq(context.TODO()), gomock.Eq(&eks.DeleteAddonInput{
						AddonName:   &addon1Name,
						ClusterName: &clusterName,
						Preserve:    true,
					})).
					Return(&eks.DeleteAddonOutput{
						Addon: &ekstypes.Addon{
							AddonArn:     aws.String(addonARN),
							AddonName:    aws.String(addon1Name),
							AddonVersion: aws.String(addon1version),
							ClusterName:  aws.String(clusterName),
							CreatedAt:    &created,
							ModifiedAt:   &created,
							Status:       ekstypes.AddonStatusDeleting,
							Tags:         createTags(),
						},
					}, nil)
				m.WaitUntilAddonDeleted(gomock.Eq(context.TODO()), gomock.Eq(&eks.DescribeAddonInput{
					AddonName:   aws.String(addon1Name),
					ClusterName: aws.String(clusterName),
				}), maxActiveUpdateDeleteWait).Return(nil)
			},
			installedAddons: []*EKSAddon{
				createInstalledAddon(addon1Name, addon1version, addonARN, addonStatusActive),
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name: "1 installed and 0 desired - addon has status of deleting",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.WaitUntilAddonDeleted(gomock.Eq(context.TODO()), gomock.Eq(&eks.DescribeAddonInput{
					AddonName:   aws.String(addon1Name),
					ClusterName: aws.String(clusterName),
				}), maxActiveUpdateDeleteWait).Return(nil)
			},
			installedAddons: []*EKSAddon{
				createInstalledAddon(addon1Name, addon1version, addonARN, addonStatusDeleting),
			},
			expectCreateError: false,
			expectDoError:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)
			tc.expect(eksMock.EXPECT())

			ctx := context.TODO()

			planner := NewPlan(clusterName, tc.desiredAddons, tc.installedAddons, eksMock, maxActiveUpdateDeleteWait, tc.preserveOnDelete)
			procedures, err := planner.Create(ctx)
			if tc.expectCreateError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
			g.Expect(procedures).NotTo(BeNil())

			for _, proc := range procedures {
				procErr := proc.Do(ctx)
				if tc.expectDoError {
					g.Expect(procErr).To(HaveOccurred())
					return
				}
				g.Expect(procErr).To(BeNil())
			}
		})
	}
}

func createTags() infrav1.Tags {
	tags := infrav1.Tags{}
	tags["tag1"] = "val1"

	return tags
}

func createTagsAdditional() infrav1.Tags {
	tags := createTags()
	tags["tag2"] = "val2"

	return tags
}

func createDesiredAddon(name, version string) *EKSAddon {
	tags := createTags()

	return &EKSAddon{
		Name:            &name,
		Version:         &version,
		Tags:            tags,
		ResolveConflict: aws.String(string(ekstypes.ResolveConflictsOverwrite)),
	}
}

func createDesiredAddonExtraTag(name, version string) *EKSAddon {
	tags := createTagsAdditional()

	return &EKSAddon{
		Name:            &name,
		Version:         &version,
		Tags:            tags,
		ResolveConflict: aws.String(string(ekstypes.ResolveConflictsOverwrite)),
	}
}

func createInstalledAddon(name, version, arn, status string) *EKSAddon {
	desired := createDesiredAddon(name, version)
	desired.ARN = &arn
	desired.Status = &status

	return desired
}
