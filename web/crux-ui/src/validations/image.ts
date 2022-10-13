import * as yup from 'yup'
import { markerRule } from './common'
import { containerConfigSchema, explicitContainerConfigSchema, uniqueKeyValuesSchema } from './container'

export const imageSchema = yup.object().shape({
  environment: uniqueKeyValuesSchema,
  annotations: markerRule,
  labels: markerRule,
  config: containerConfigSchema,
})

export const patchContainerConfigSchema = yup.object().shape({
  environment: uniqueKeyValuesSchema.nullable(),
  capabilities: uniqueKeyValuesSchema.nullable(),
  annotations: markerRule,
  labels: markerRule,
  config: explicitContainerConfigSchema.nullable(),
})
