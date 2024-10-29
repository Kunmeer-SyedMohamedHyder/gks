// Package greenscheduling provides a Kubernetes scheduling plugin designed to prioritize node selection
// based on environmental sustainability factors. By integrating with the SIC (Sustainability Insight
// Center) API, this plugin gathers data on carbon emissions, energy usage, and other sustainability metrics,
// using them to calculate a sustainability score for each node. This score can guide Kubernetes in selecting
// nodes that contribute to more environmentally-friendly resource allocation decisions.
//
// The GreenScheduling plugin is intended for organizations and platforms aiming to enhance the environmental
// sustainability of their cloud infrastructure by reducing carbon footprints, optimizing energy usage,
// and potentially lowering costs. By leveraging this plugin, Kubernetes clusters can integrate sustainability
// as a core consideration within their scheduling operations.
//
// # Example Usage
//
// To use the GreenScheduling plugin within a Kubernetes environment, create a new instance of the plugin
// by providing the necessary configuration, Kubernetes client, and SIC client. Configure `TimeSeriesConfig`
// and `SustainabilityWeights` based on organizational priorities for sustainability. Once integrated with
// the scheduler, Kubernetes will use the plugin’s calculated scores to select nodes, thereby optimizing
// for environmental impact.
//
// # Future Enhancements
//
// Future versions of the `greenscheduling` package may include additional sustainability metrics,
// more advanced scoring models, and deeper integrations with other environmental data sources. Plans also
// include providing more configuration options for end-users to further refine how each metric contributes
// to the sustainability score, aligning the plugin with evolving environmental goals and regulatory standards.
package greenscheduling
