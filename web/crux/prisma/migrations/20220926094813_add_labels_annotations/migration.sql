-- AlterTable
ALTER TABLE "ContainerConfig" ADD COLUMN     "annotations" JSONB NOT NULL DEFAULT '{}',
ADD COLUMN     "labels" JSONB NOT NULL DEFAULT '{}';

-- AlterTable
ALTER TABLE "InstanceContainerConfig" ADD COLUMN     "annotations" JSONB,
ADD COLUMN     "labels" JSONB;
