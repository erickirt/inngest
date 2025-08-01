import { Select, type Option } from '../Select/Select';
import { FunctionRunTimeField, isFunctionTimeField } from '../types/functionRun';

type TimeFieldFilterProps = {
  selectedTimeField: FunctionRunTimeField;
  onTimeFieldChange: (value: FunctionRunTimeField) => void;
};

const options: Option[] = [
  FunctionRunTimeField.QueuedAt,
  FunctionRunTimeField.StartedAt,
  FunctionRunTimeField.EndedAt,
].map((field) => ({
  id: field,
  name: field.replace(/_/g, ' '),
}));

export default function TimeFieldFilter({
  selectedTimeField,
  onTimeFieldChange,
}: TimeFieldFilterProps) {
  const selectedValue = options.find((option) => option.id === selectedTimeField);

  return (
    <Select
      value={selectedValue}
      onChange={(value: Option) => {
        if (isFunctionTimeField(value.id)) {
          onTimeFieldChange(value.id);
        }
      }}
      label="Time Field"
      isLabelVisible={false}
      className="bg-modalBase min-w-[90px]"
      size="small"
    >
      <Select.Button size="small">
        <span className="lowercase first-letter:capitalize">{selectedValue?.name}</span>
      </Select.Button>
      <Select.Options>
        {options.map((option) => {
          return (
            <Select.Option key={option.id} option={option}>
              <span className="inline-flex items-center gap-2 lowercase">
                <label className="first-letter:capitalize">{option.name}</label>
              </span>
            </Select.Option>
          );
        })}
      </Select.Options>
    </Select>
  );
}
