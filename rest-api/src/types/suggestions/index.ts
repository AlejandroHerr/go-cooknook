// Code generated by tygo. DO NOT EDIT.

//////////
// source: model.go

export interface Option {
  label: string;
  value: string;
}

//////////
// source: router.go

export interface GetSuggestionsReponse {
  options: Option[];
}
