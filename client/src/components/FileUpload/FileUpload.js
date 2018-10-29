import React, { Component } from 'react';
import DropzoneComponent from 'react-dropzone-component';
import { socketSuccess, socketError, getSocket } from './api';
import {NotificationContainer, NotificationManager} from 'react-notifications';
import axios from 'axios';

import 'react-notifications/lib/notifications.css';
import './FileUpload.css'

class FileUpload extends Component {
    constructor(props) {
        super(props);

        this.state = {
            statusMsg: '',
            statusCode: '',
            socketMsg: '',
            socketId: ''
        };

        this.djsConfig = {
            addRemoveLinks: true,
            acceptedFiles: "video/mp4",
            autoProcessQueue: false,
            autoDiscover: false,
            maxFilesize: 2048
        };

        this.componentConfig = {
            iconFiletypes: ['.jpg', '.png', '.gif'],
            showFiletypeIcon: true,
            postUrl: 'localhost:8888/upload'
        };

        // socket.io
        socketSuccess((err, socketMsg) =>
            NotificationManager.success(socketMsg, 'Success')
        );

        socketError((err, socketMsg) =>
            NotificationManager.error(socketMsg, 'Error')
        );

    }

    componentDidMount() {
        const socket = getSocket();

        socket.on('connect', () => {
            this.setState({socketId:  socket.id})
        });

    }

    handleFileSubmit() {
        this.dropzone.getAcceptedFiles().map(file => {
            axios.put('http://localhost:8888/upload?callback=http://localhost:3000',file, {
                headers: { socketId: this.state.socketId }
            })
                .then(response => {
                    this.setState({statusCode: response.status, statusMsg: response.data});
                    file.previewElement.classList.remove("dz-error");
                    file.previewElement.classList.add("dz-success");
                    NotificationManager.success('File uploaded', 'Success');
                })
                .catch(error => {
                    if (!! error.response){
                        this.setState({statusCode: error.response.status, statusMsg: error.response.data});
                        file.previewElement.classList.remove("dz-success");
                        file.previewElement.classList.add("dz-error");
                        file.previewElement.children[3].innerText = error.response.data;
                        NotificationManager.error(error.response.data, 'Error');
                    } else {
                        const status = 404;
                        const data = "Server not found";
                        this.setState({statusCode: status, statusMsg: data});
                        file.previewElement.classList.remove("dz-success");
                        file.previewElement.classList.add("dz-error");
                        file.previewElement.children[3].innerText = data;
                        NotificationManager.error(data, 'Error');
                    }

                });
        })
    }

    render() {
        const config = this.componentConfig;
        const djsConfig = this.djsConfig;

        const eventHandlers = {
            init: dz => this.dropzone = dz,
        };

        return (
            <div>
                <DropzoneComponent config={config} eventHandlers={eventHandlers} djsConfig={djsConfig} />
                <button id="submit-all" onClick={this.handleFileSubmit.bind(this)}>Submit all files</button>
                <NotificationContainer/>
            </div>
        );
    }
}

export default FileUpload;