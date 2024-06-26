import React, { useLayoutEffect, useMemo } from 'react';
import * as am5 from "@amcharts/amcharts5";
import * as am5hierarchy from "@amcharts/amcharts5/hierarchy";
import am5themes_Animated from "@amcharts/amcharts5/themes/Animated";

const genders = { male: "#5af", female: "#f66" }

const proccess = data => ({
    ...data.animal,
    color: genders[data.animal.gender] || "#666",
    children: (data?.children || []).map(proccess),
})

const AnimalTree = ({ data: dataInput = {} }) => {
    const divID = useMemo(() => ("" + Math.random()).substr(2, 10), [])
    useLayoutEffect(() => {
        let root = am5.Root.new(divID);

        root.setThemes([
            am5themes_Animated.new(root)
        ]);

        // Chart data
        const data = proccess(dataInput)

        var zoomableContainer = root.container.children.push(
            am5.ZoomableContainer.new(root, {
                width: am5.p100,
                height: am5.p100,
                wheelable: true,
                pinchZoom: true
            })
        );

        var zoomTools = zoomableContainer.children.push(am5.ZoomTools.new(root, {
            target: zoomableContainer
        }));

        // Create series
        // https://www.amcharts.com/docs/v5/charts/hierarchy/#Adding
        var series = zoomableContainer.contents.children.push(am5hierarchy.Tree.new(root, {
            singleBranchOnly: false,
            downDepth: 1,
            initialDepth: 10,
            valueField: "value",
            categoryField: "name",
            childDataField: "children",
            fillField: "color",
            inversed: true
        }));

        series.labels.template.set("minScale", 0);

        series.data.setAll([data]);
        series.set("selectedDataItem", series.dataItems[0]);

        // Chart stop

        return () => {
            root.dispose();
        };
    }, [dataInput]);

    return (
        <div id={divID} style={{ width: "100%", height: "100%" }}></div>
    );
}

export default AnimalTree