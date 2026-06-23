import { customizeValidator } from '@rjsf/validator-ajv8'
import Ajv2020 from 'ajv/dist/2020'

const validator = customizeValidator({ AjvClass: Ajv2020 as any })

export default validator
