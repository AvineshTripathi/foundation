// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the EARLY ACCESS SOFTWARE LICENSE AGREEMENT
// available at http://github.com/namespacelabs/foundation

import classNames from "classnames";
import { Link } from "wouter";
import {
	DataType,
	ForwardedPort,
	ServerType,
	StackEntryStateType,
	StackType,
} from "../../datamodel/Schema";
import { useTasksByServer } from "../../datamodel/TasksObserver";
import { ExternalLinkIcon } from "../../icons";
import Selectable from "../../ui/sidebar/Selectable";
import Tabs from "../../ui/sidebar/Tabs";
import { useServerRoute } from "../server/routing";
import classes from "./sidebar.module.css";

export default function ServerBlock(props: { data: DataType }) {
	let [matches, params] = useServerRoute();

	let tabs = [
		{
			id: "stack",
			label: "Stack",
			render: () => (
				<>
					{props.data.stack?.entry?.map((s) => {
						return (
							<Selectable
								key={s.server.package_name}
								selected={matches && params?.id === s.server.id}>
								<Server server={s.server} state={stateOf(props.data, s.server.package_name)} />
							</Selectable>
						);
					})}
				</>
			),
		},
	];

	return (
		<div className={classes.serverContent}>
			{/* <ServerSelector data={props.data} /> */}

			<Tabs tabs={tabs} />

			<ForwardedPorts data={props.data} />
		</div>
	);
}

function stateOf(data: DataType, packageName: string) {
	let matchingState = data.state?.filter((st) => st.package_name === packageName);
	return matchingState?.length
		? matchingState[0]
		: ({ package_name: packageName } as StackEntryStateType);
}

const taskHumanNames: { [key: string]: string } = {
	"graph.compute": "computing",
	"server.build": "building",
	"server.provision": "provisioning",
	"server.deploy": "deploying",
	"server.start": "starting",
};

function humanTaskName(name: string) {
	return taskHumanNames[name] || name;
}

function Server(props: { server: ServerType; state: StackEntryStateType; stack?: StackType }) {
	let runningTask = useTasksByServer(props.server.package_name);

	const parts = props.server.package_name.split("/");
	let p = parts.pop();

	while (parts.length) {
		let would = parts.pop() + "/" + p;
		if (would.length > 34) {
			break;
		}
		p = would;
	}

	if (parts.length) {
		p = "... " + p;
	}

	let badges: string[] = [];

	let isWorking = false;
	if (props.state.last_error) {
		p = "failed: " + props.state.last_error;
	} else if (runningTask.length) {
		// Show the last task, and collapse the rest into "...".
		if (runningTask.length > 1) {
			badges.push("...");
		}

		badges.push(humanTaskName(runningTask[runningTask.length - 1].name));
	}

	return (
		<Link href={`/server/${props.server.id}`}>
			<a className={classes.serverItem}>
				<div className={classes.serverName}>
					<span>{props.server.name}</span>
					{!isWorking &&
						badges.map((b) => (
							<span key={b} className={classes.badge}>
								{b}
							</span>
						))}
				</div>
				<div className={classes.serverPackageName}>
					{isWorking ? (
						badges.map((b) => (
							<span key={b} className={classes.badge}>
								{b}
							</span>
						))
					) : (
						<span>{p}</span>
					)}
				</div>
			</a>
		</Link>
	);
}

function ForwardedPorts(props: { data: DataType }) {
	if (!props.data.forwarded_port) {
		return null;
	}

	let tabs = [
		{
			id: "ports",
			label: "Ports",
			render: () => (
				<>
					{sortPorts(props.data.current.server.package_name, props.data.forwarded_port).map((p) => (
						<Port key={p.container_port} data={props.data} p={p} />
					))}
				</>
			),
		},
	];

	return <Tabs tabs={tabs} />;
}

function sortPorts(current: string, ports?: ForwardedPort[]) {
	let copy = [...(ports || [])];

	copy.sort((a, b) => {
		let apkg = a.endpoint.server_owner || "<stack>";
		let bpkg = b.endpoint.server_owner || "<stack>";

		if (apkg === current) {
			return -1;
		} else if (bpkg === current) {
			return 1;
		}

		if (apkg === bpkg) {
			return a.container_port - b.container_port;
		}
		return apkg.localeCompare(bpkg);
	});

	return copy;
}

function Port(props: { data: DataType; p: ForwardedPort }) {
	let { p } = props;
	let serverName = p.endpoint.server_owner;

	for (let s of props.data.stack?.entry || []) {
		if (s.server.package_name === p.endpoint.server_owner) {
			serverName = s.server.name;
			break;
		}
	}

	let body = (
		<>
			<div className={classes.serviceDesc}>
				<span className={classes.ports}>{p.local_port}</span>{" "}
				<span className={classes.serviceName}>{serverName} </span>
				<ExternalLinkIcon />
			</div>
			<div className={classes.serviceDetails}>
				{p.endpoint.port?.name || p.endpoint.service_name} ({p.container_port})
			</div>
		</>
	);

	return (
		<a
			className={classes.port}
			href={`${window.location.protocol}//${window.location.hostname}:${p.local_port}`}
			target="_blank">
			{body}
		</a>
	);
}
