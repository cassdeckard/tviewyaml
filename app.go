package tviewyaml

// This file previously contained the CreateApp function.
// The API has been replaced with the Builder pattern.
// Use NewAppBuilder(configDir).Build() instead.
//
// Example:
//   app, err := tviewyaml.NewAppBuilder("./config").Build()
//
// For custom template functions:
//   app, err := tviewyaml.NewAppBuilder("./config").
//       WithTemplateFunction("myFunc", 1, intPtr(1), nil, handler).
//       Build()

