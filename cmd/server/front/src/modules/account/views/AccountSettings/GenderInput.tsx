import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import React from 'react'
import { z } from 'zod'

export type GenderInputProps = {
  value: string
  onChange: (value: string) => void
}

const GENDER_MALE = 'Male'
const GENDER_FEMALE = 'Female'
const GENDER_NONE = ''

export default function GenderInput({ value, onChange }: GenderInputProps) {
  const [genderType, setGenderType] = React.useState(getGenderType(value))

  return (
    <div className="space-y-4">
      <RadioGroup
        value={genderType}
        onValueChange={(str) => {
          const value = genderTypeSchema.parse(str)
          setGenderType(value)
          if (value === 'f') onChange(GENDER_FEMALE)
          else if (value === 'm') onChange(GENDER_MALE)
          else if (value === 'none') onChange(GENDER_NONE)
        }}
      >
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="none" id="gender-type-none" />
          <Label htmlFor="gender-type-none">Rather not say</Label>
        </div>
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="m" id="gender-type-m" />
          <Label htmlFor="gender-type-m">Male</Label>
        </div>
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="f" id="gender-type-f" />
          <Label htmlFor="gender-type-f">Female</Label>
        </div>
        <div className="flex items-center space-x-2">
          <RadioGroupItem value="custom" id="gender-type-custom" />
          <Label htmlFor="gender-type-custom">Custom</Label>
        </div>
      </RadioGroup>
      {genderType === 'custom' && (
        <Input value={value} onChange={(e) => onChange(e.target.value)} />
      )}
    </div>
  )
}

const genderTypeSchema = z.enum(['m', 'f', 'none', 'custom'])

type GenderType = z.infer<typeof genderTypeSchema>

function getGenderType(value: string): GenderType {
  if (value === GENDER_MALE) return 'm'
  if (value === GENDER_FEMALE) return 'f'
  if (value === GENDER_NONE) return 'none'
  return 'custom'
}
