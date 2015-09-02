var WorkflowBox = React.createClass({
  loadWorkflowsFromServer: function() {
    $.ajax({
      url: this.props.url,
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({data: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },

  getInitialState: function() {
    return {data: []};
  },

  componentDidMount: function() {
    this.loadWorkflowsFromServer();
    setInterval(this.loadWorkflowsFromServer, this.props.pollInterval);
  },

  render: function() {
    return (
      <div className="workflowBox dropdown">
	<button className="mdl-button mdl-js-button mdl-button--raised mdl-js-ripple-effect" type="button" data-toggle="dropdown">Workflows 
  	</button> 
	<ul className="dropdown-menu">
        	<WorkflowList data={this.state.data} />
	</ul>
      </div>
    );
  }
});


var WorkflowList = React.createClass({
  render: function() {
    var workflowNodes = this.props.data.map(function (workflow) {
      return (
        <Workflow name={workflow.name} id={workflow.id}>
        </Workflow>
      );
    });
    return (
      <div className="workflowList">
        {workflowNodes}
      </div>
    );
  }
});

var Workflow = React.createClass({
    getInitialState: function() {
        return {
            isSelected: false
        };
    },
    handleClick: function() {
        this.setState({
            isSelected: true
        })
	React.render(
	<CommentForm name={this.props.name} id={this.props.id} />,
	document.getElementById('workflow-form')	
);
    },
  render: function() {
    return (
      <div className="workflow">
        <li onClick={this.handleClick} className="workflowName"><a href="#">
	<WFName name={this.props.name} />
        </a></li>
        {this.props.children}
      </div>
    );
  }
});


var WFName = React.createClass({
  render: function() {
    return (
<div>{this.props.name}</div>
    );
  }
});

var CommentForm = React.createClass({
  render: function() {
    return (
	<div>{this.props.id}</div>
    );
  }
});

React.render(
<WorkflowBox url="/v1/workflows/all" pollInterval={2000} />,
  document.getElementById('workflow-dropdown')
);

