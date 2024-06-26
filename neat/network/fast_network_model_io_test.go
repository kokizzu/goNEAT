package network

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const jsonFMNStr = `{"id":123456,"name":"test network","input_neuron_count":2,"sensor_neuron_count":3,"output_neuron_count":2,"bias_neuron_count":1,"total_neuron_count":8,"activation_functions":["SigmoidSteepenedActivation","SigmoidSteepenedActivation","SigmoidSteepenedActivation","SigmoidSteepenedActivation","SigmoidSteepenedActivation","SigmoidSteepenedActivation","SigmoidSteepenedActivation","SigmoidSteepenedActivation"],"bias_list":[0,0,0,0,0,0,1,0],"connections":[{"source_index":1,"target_index":5,"weight":15,"signal":0},{"source_index":2,"target_index":5,"weight":10,"signal":0},{"source_index":2,"target_index":6,"weight":5,"signal":0},{"source_index":6,"target_index":7,"weight":17,"signal":0},{"source_index":5,"target_index":3,"weight":7,"signal":0},{"source_index":7,"target_index":3,"weight":4.5,"signal":0},{"source_index":7,"target_index":4,"weight":13,"signal":0}]}`
const jsonFNMStrModule = `{"id":123456,"name":"test network","input_neuron_count":2,"sensor_neuron_count":3,"output_neuron_count":2,"bias_neuron_count":1,"total_neuron_count":8,"activation_functions":["SigmoidSteepenedActivation","SigmoidSteepenedActivation","SigmoidSteepenedActivation","LinearActivation","LinearActivation","LinearActivation","LinearActivation","NullActivation"],"bias_list":[0,0,0,0,0,10,1,0],"connections":[{"source_index":1,"target_index":5,"weight":15,"signal":0},{"source_index":2,"target_index":6,"weight":5,"signal":0},{"source_index":7,"target_index":3,"weight":4.5,"signal":0},{"source_index":7,"target_index":4,"weight":13,"signal":0}],"modules":[{"activation_type":"MultiplyModuleActivation","input_indexes":[5,6],"output_indexes":[7]}]}`

const networkName = "test network"
const networkId = 123456

func TestFastModularNetworkSolver_WriteModel_NoModule(t *testing.T) {
	net := buildNamedNetwork(networkName, networkId)

	fmm, err := net.FastNetworkSolver()
	require.NoError(t, err, "failed to create fast network solver")

	outBuf := bytes.NewBufferString("")
	err = fmm.(*FastModularNetworkSolver).WriteModel(outBuf)
	require.NoError(t, err, "failed to write model")

	println(outBuf.String())

	var expected interface{}
	err = json.Unmarshal([]byte(jsonFMNStr), &expected)
	require.NoError(t, err, "failed to unmarshal expected json")
	var actual interface{}
	err = json.Unmarshal(outBuf.Bytes(), &actual)
	require.NoError(t, err, "failed to unmarshal actual json")

	assert.EqualValues(t, expected, actual, "model JSON does not match expected JSON")
}

func TestFastModularNetworkSolver_WriteModel_WithModule(t *testing.T) {
	net := buildNamedModularNetwork(networkName, networkId)

	fmm, err := net.FastNetworkSolver()
	require.NoError(t, err, "failed to create fast network solver")

	outBuf := bytes.NewBufferString("")
	err = fmm.(*FastModularNetworkSolver).WriteModel(outBuf)
	require.NoError(t, err, "failed to write model")

	println(outBuf.String())

	var expected interface{}
	err = json.Unmarshal([]byte(jsonFNMStrModule), &expected)
	require.NoError(t, err, "failed to unmarshal expected json")
	var actual interface{}
	err = json.Unmarshal(outBuf.Bytes(), &actual)
	require.NoError(t, err, "failed to unmarshal actual json")

	assert.EqualValues(t, expected, actual, "model JSON does not match expected JSON")
}

func TestReadFMNSModel_NoModule(t *testing.T) {
	buf := bytes.NewBufferString(jsonFMNStr)

	fmm, err := ReadFMNSModel(buf)
	assert.NoError(t, err, "failed to read model")
	assert.NotNil(t, fmm, "failed to deserialize model")

	assert.Equal(t, fmm.Name, networkName, "wrong network name")
	assert.Equal(t, fmm.Id, networkId, "wrong network id")

	data := []float64{1.5, 2.0} // bias inherent
	err = fmm.LoadSensors(data)
	require.NoError(t, err, "failed to load sensors")

	// test that it operates as expected
	//
	net := buildNetwork()
	depth, err := net.MaxActivationDepth()
	require.NoError(t, err, "failed to calculate max depth")

	t.Logf("depth: %d\n", depth)
	logNetworkActivationPath(net, t)

	data = append(data, 1.0) // BIAS is third object
	err = net.LoadSensors(data)
	require.NoError(t, err, "failed to load sensors")
	res, err := net.ForwardSteps(depth)
	require.NoError(t, err, "error when trying to activate objective network")
	require.True(t, res, "failed to activate objective network")

	// do forward steps through the solver and test results
	//
	res, err = fmm.Relax(depth, .1)
	require.NoError(t, err, "error while do forward steps")
	require.True(t, res, "forward steps returned false")

	// check results by comparing activations of objective network and fast network solver
	//
	outputs := fmm.ReadOutputs()
	for i, out := range outputs {
		assert.Equal(t, net.Outputs[i].Activation, out, "wrong activation at: %d", i)
	}
}

func TestReadFMNSModel_ModularNetwork(t *testing.T) {
	buf := bytes.NewBufferString(jsonFNMStrModule)

	fmm, err := ReadFMNSModel(buf)
	assert.NoError(t, err, "failed to read model")
	assert.NotNil(t, fmm, "failed to deserialize model")

	assert.Equal(t, fmm.Name, networkName, "wrong network name")
	assert.Equal(t, fmm.Id, networkId, "wrong network id")

	data := []float64{1.0, 2.0} // bias inherent
	err = fmm.LoadSensors(data)
	require.NoError(t, err, "failed to load sensors")

	// test that it operates as expected
	//
	net := buildModularNetwork()
	depth, err := net.MaxActivationDepth()
	require.NoError(t, err, "failed to calculate max depth")

	t.Logf("depth: %d\n", depth)
	logNetworkActivationPath(net, t)

	// activate objective network
	//
	data = append(data, 1.0) // BIAS is third object
	err = net.LoadSensors(data)
	require.NoError(t, err, "failed to load sensors")
	res, err := net.ForwardSteps(depth)
	require.NoError(t, err, "error when trying to activate objective network")
	require.True(t, res, "failed to activate objective network")

	// do forward steps through the solver and test results
	//
	res, err = fmm.Relax(depth, 1)
	require.NoError(t, err, "error while do forward steps")
	require.True(t, res, "forward steps returned false")

	// check results by comparing activations of objective network and fast network solver
	//
	outputs := fmm.ReadOutputs()
	for i, out := range outputs {
		assert.Equal(t, net.Outputs[i].Activation, out, "wrong activation at: %d", i)
	}

}
func buildNamedNetwork(name string, id int) *Network {
	net := buildNetwork()
	net.Name = name
	net.Id = id
	return net
}

func buildNamedModularNetwork(name string, id int) *Network {
	net := buildModularNetwork()
	net.Name = name
	net.Id = id
	return net
}
