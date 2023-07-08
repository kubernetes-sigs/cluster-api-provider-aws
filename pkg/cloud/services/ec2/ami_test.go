/*
Copyright 2019 The Kubernetes Authors.

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

package ec2

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ssm/mock_ssmiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
)

func TestDefaultAMILookup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		ownerID           string
		baseOS            string
		architecture      string
		kubernetesVersion string
		amiNameFormat     string
	}

	testCases := []struct {
		name   string
		args   args
		expect func(m *mocks.MockEC2APIMockRecorder)
		check  func(g *WithT, img *ec2.Image, err error)
	}{
		{
			name: "Should return latest AMI in case of valid inputs",
			args: args{
				ownerID:           "ownerID",
				baseOS:            "baseOS",
				architecture:      "x86_64",
				kubernetesVersion: "v1.0.0",
				amiNameFormat:     "ami-name",
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								ImageId:      aws.String("ancient"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("latest"),
								CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("oldest"),
								CreationDate: aws.String("2014-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
			},
			check: func(g *WithT, img *ec2.Image, err error) {
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(*img.ImageId).Should(ContainSubstring("latest"))
			},
		},
		{
			name: "Should return with error if AWS DescribeImages call failed with some error",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(nil, awserrors.NewFailedDependency("dependency failure"))
			},
			check: func(g *WithT, img *ec2.Image, err error) {
				g.Expect(err).To(HaveOccurred())
				g.Expect(img).To(BeNil())
			},
		},
		{
			name: "Should return with error if empty list of images returned from AWS ",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{}, nil)
			},
			check: func(g *WithT, img *ec2.Image, err error) {
				g.Expect(err).To(HaveOccurred())
				g.Expect(img).To(BeNil())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			ec2Mock := mocks.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock.EXPECT())

			img, err := DefaultAMILookup(ec2Mock, tc.args.ownerID, tc.args.baseOS, tc.args.kubernetesVersion, tc.args.architecture, tc.args.amiNameFormat)
			tc.check(g, img, err)
		})
	}
}

func TestDefaultAMILookupArm64(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		ownerID           string
		baseOS            string
		architecture      string
		kubernetesVersion string
		amiNameFormat     string
	}

	testCases := []struct {
		name   string
		args   args
		expect func(m *mocks.MockEC2APIMockRecorder)
		check  func(g *WithT, img *ec2.Image, err error)
	}{
		{
			name: "Should return latest AMI in case of valid inputs",
			args: args{
				ownerID:           "ownerID",
				baseOS:            "baseOS",
				architecture:      "arm64",
				kubernetesVersion: "v1.0.0",
				amiNameFormat:     "ami-name",
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								ImageId:      aws.String("ancient"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("latest"),
								CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("oldest"),
								CreationDate: aws.String("2014-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
			},
			check: func(g *WithT, img *ec2.Image, err error) {
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(*img.ImageId).Should(ContainSubstring("latest"))
			},
		},
		{
			name: "Should return with error if AWS DescribeImages call failed with some error",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(nil, awserrors.NewFailedDependency("dependency failure"))
			},
			check: func(g *WithT, img *ec2.Image, err error) {
				g.Expect(err).To(HaveOccurred())
				g.Expect(img).To(BeNil())
			},
		},
		{
			name: "Should return with error if empty list of images returned from AWS ",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{}, nil)
			},
			check: func(g *WithT, img *ec2.Image, err error) {
				g.Expect(err).To(HaveOccurred())
				g.Expect(img).To(BeNil())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			ec2Mock := mocks.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock.EXPECT())

			img, err := DefaultAMILookup(ec2Mock, tc.args.ownerID, tc.args.baseOS, tc.args.kubernetesVersion, tc.args.architecture, tc.args.amiNameFormat)
			tc.check(g, img, err)
		})
	}
}
func TestAMIs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		expect func(m *mocks.MockEC2APIMockRecorder)
		check  func(g *WithT, id string, err error)
	}{
		{
			name: "Should return latest AMI in case of valid inputs",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								ImageId:      aws.String("ancient"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("latest"),
								CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("oldest"),
								CreationDate: aws.String("2014-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
			},
			check: func(g *WithT, id string, err error) {
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(id).Should(ContainSubstring("latest"))
			},
		},
		{
			name: "Should return error if invalid creation date passed",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								ImageId:      aws.String("ancient"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("latest"),
								CreationDate: aws.String("invalid creation date"),
							},
							{
								ImageId:      aws.String("oldest"),
								CreationDate: aws.String("2014-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
			},
			check: func(g *WithT, id string, err error) {
				g.Expect(err).To(HaveOccurred())
				g.Expect(id).Should(BeEmpty())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			scheme, err := setupScheme()
			g.Expect(err).NotTo(HaveOccurred())
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			ec2Mock := mocks.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock.EXPECT())

			clusterScope, err := setupClusterScope(client)
			g.Expect(err).NotTo(HaveOccurred())

			s := NewService(clusterScope)
			s.EC2Client = ec2Mock

			id, err := s.defaultAMIIDLookup("", "", "base os-baseos version", "x86_64", "v1.11.1")
			tc.check(g, id, err)
		})
	}
}

func TestFormatVersionForEKS(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
		wantErr bool
	}{
		{
			name:    "Should remove non zero patch from version",
			version: "v1.23.2",
			want:    "1.23",
			wantErr: false,
		},
		{
			name:    "Should return major.minor in case patch is nil",
			version: "v1.23",
			want:    "1.23",
			wantErr: false,
		},
		{
			name:    "Should return minor as zero if only major is present in version",
			version: "v1",
			want:    "1.0",
			wantErr: false,
		},
		{
			name:    "Should return error if invalid version is given",
			version: "v1-23.3",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			got, err := formatVersionForEKS(tt.version)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(got).Should(BeEquivalentTo(tt.want))
		})
	}
}

func TestGenerateAMIName(t *testing.T) {
	type args struct {
		amiNameFormat     string
		baseOS            string
		kubernetesVersion string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should return image name even if OS and amiNameFormat is empty",
			args: args{
				kubernetesVersion: "1.23.3",
			},
			want: "capa-ami--?1.23.3-*",
		},
		{
			name: "Should return valid amiName if default AMI name format passed",
			args: args{
				amiNameFormat:     DefaultAmiNameFormat,
				baseOS:            "centos-7",
				kubernetesVersion: "1.23.3",
			},
			want: "capa-ami-centos-7-?1.23.3-*",
		},
		{
			name: "Should return valid amiName if custom AMI name format passed",
			args: args{
				amiNameFormat:     "random-{{.BaseOS}}-?{{.K8sVersion}}-*",
				baseOS:            "centos-7",
				kubernetesVersion: "1.23.3",
			},
			want: "random-centos-7-?1.23.3-*",
		},
		{
			name: "Should return valid amiName if new AMI name format passed",
			args: args{
				amiNameFormat:     "random-{{.BaseOS}}-{{.K8sVersion}}",
				baseOS:            "centos-7",
				kubernetesVersion: "v1.23.3",
			},
			want: "random-centos-7-v1.23.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			got, err := GenerateAmiName(tt.args.amiNameFormat, tt.args.baseOS, tt.args.kubernetesVersion)
			g.Expect(err).To(BeNil())
			g.Expect(got).Should(Equal(tt.want))
		})
	}
}

func TestGetLatestImage(t *testing.T) {
	tests := []struct {
		name    string
		imgs    []*ec2.Image
		want    *ec2.Image
		wantErr bool
	}{
		{
			name: "Should return image with latest creation date",
			imgs: []*ec2.Image{
				{
					ImageId:      aws.String("ancient"),
					CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
				},
				{
					ImageId:      aws.String("latest"),
					CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
				},
				{
					ImageId:      aws.String("oldest"),
					CreationDate: aws.String("2014-02-08T17:02:31.000Z"),
				},
			},
			want: &ec2.Image{
				ImageId:      aws.String("latest"),
				CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
			},
			wantErr: false,
		},
		{
			name: "Should return last image if all images have same creation date",
			imgs: []*ec2.Image{
				{
					ImageId:      aws.String("image 1"),
					CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
				},
				{
					ImageId:      aws.String("image 2"),
					CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
				},
				{
					ImageId:      aws.String("image 3"),
					CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
				},
			},
			want: &ec2.Image{
				ImageId:      aws.String("image 3"),
				CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
			},
			wantErr: false,
		},
		{
			name: "Should return error if creation date is given in wrong format",
			imgs: []*ec2.Image{
				{
					ImageId:      aws.String("image 1"),
					CreationDate: aws.String("2019-02-08"),
				},
				{
					ImageId:      aws.String("image 2"),
					CreationDate: aws.String("2019-02-08"),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			got, err := GetLatestImage(tt.imgs)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(got).Should(Equal(tt.want))
		})
	}
}

func TestEKSAMILookUp(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	gpuAMI := infrav1.AmazonLinuxGPU
	tests := []struct {
		name       string
		k8sVersion string
		arch       string
		amiType    *infrav1.EKSAMILookupType
		expect     func(m *mock_ssmiface.MockSSMAPIMockRecorder)
		want       string
		wantErr    bool
	}{
		{
			name:       "Should return an id corresponding to GPU if GPU based AMI type passed",
			k8sVersion: "v1.23.3",
			arch:       "x86_64",
			amiType:    &gpuAMI,
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.GetParameter(gomock.Eq(&ssm.GetParameterInput{
					Name: aws.String("/aws/service/eks/optimized-ami/1.23/amazon-linux-2-gpu/recommended/image_id"),
				})).Return(&ssm.GetParameterOutput{
					Parameter: &ssm.Parameter{
						Value: aws.String("id"),
					},
				}, nil)
			},
			want:    "id",
			wantErr: false,
		},
		{
			name:       "Should return an id not corresponding to GPU if AMI type is default",
			k8sVersion: "v1.23.3",
			arch:       "x86_64",
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.GetParameter(gomock.Eq(&ssm.GetParameterInput{
					Name: aws.String("/aws/service/eks/optimized-ami/1.23/amazon-linux-2/recommended/image_id"),
				})).Return(&ssm.GetParameterOutput{
					Parameter: &ssm.Parameter{
						Value: aws.String("id"),
					},
				}, nil)
			},
			want:    "id",
			wantErr: false,
		},
		{
			name:       "Should return an error if GetParameter call fails with some AWS error",
			k8sVersion: "v1.23.3",
			arch:       "x86_64",
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.GetParameter(gomock.Eq(&ssm.GetParameterInput{
					Name: aws.String("/aws/service/eks/optimized-ami/1.23/amazon-linux-2/recommended/image_id"),
				})).Return(nil, awserrors.NewFailedDependency("dependency failure"))
			},
			wantErr: true,
		},
		{
			name:       "Should return an error if invalid Kubernetes version passed",
			k8sVersion: "__$__",
			arch:       "x86_64",
			wantErr:    true,
		},
		{
			name:       "Should return an error if no SSM parameter found",
			k8sVersion: "v1.23.3",
			arch:       "x86_64",
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.GetParameter(gomock.Eq(&ssm.GetParameterInput{
					Name: aws.String("/aws/service/eks/optimized-ami/1.23/amazon-linux-2/recommended/image_id"),
				})).Return(&ssm.GetParameterOutput{}, nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			scheme, err := setupScheme()
			g.Expect(err).NotTo(HaveOccurred())
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			ssmMock := mock_ssmiface.NewMockSSMAPI(mockCtrl)
			if tt.expect != nil {
				tt.expect(ssmMock.EXPECT())
			}

			clusterScope, err := setupClusterScope(client)
			g.Expect(err).NotTo(HaveOccurred())

			s := NewService(clusterScope)
			s.SSMClient = ssmMock

			got, err := s.eksAMILookup(tt.k8sVersion, tt.arch, tt.amiType)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(got).Should(Equal(tt.want))
		})
	}
}
