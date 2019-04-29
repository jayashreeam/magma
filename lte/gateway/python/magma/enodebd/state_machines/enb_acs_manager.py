"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
from typing import Any, Optional, List
from magma.common.service import MagmaService
from magma.enodebd.devices.device_map import get_device_handler_from_name
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine
from magma.enodebd.state_machines.acs_state_utils import \
    get_device_name_from_inform
from magma.enodebd.tr069 import models
from spyne import ComplexModelBase
from spyne.server.wsgi import WsgiMethodContext


class StateMachineManager:
    """
    Delegates tr069 message handling to a dedicated state machine for the
    device.
    """
    def __init__(
        self,
        service: MagmaService,
    ):
        self._ip_serial_mapping = IpToSerialMapping()
        self._service = service
        self._state_machine_by_ip = {}

    def handle_tr069_message(
        self,
        ctx: WsgiMethodContext,
        tr069_message: ComplexModelBase,
    ) -> Any:
        """ Delegate message handling to the appropriate eNB state machine """
        client_ip = self._get_client_ip(ctx)
        if isinstance(tr069_message, models.Inform):
            self._update_device_mapping(client_ip, tr069_message)
        handler = self._get_handler(client_ip)

        if handler is None:
            logging.warning('Received non-Inform tr069 msg from unknown eNB. '
                            'Ending session with empty HTTP response.')
            return models.DummyInput()

        return handler.handle_tr069_message(tr069_message)

    def get_handler_by_ip(self, client_ip: str) -> EnodebAcsStateMachine:
        return self._state_machine_by_ip[client_ip]

    def get_handler_by_serial(self, enb_serial: str) -> EnodebAcsStateMachine:
        client_ip = self._ip_serial_mapping.get_ip(enb_serial)
        return self._state_machine_by_ip[client_ip]

    def get_connected_serial_id_list(self) -> List[str]:
        return self._ip_serial_mapping.get_serial_list()

    def get_ip_of_serial(self, enb_serial: str) -> str:
        return self._ip_serial_mapping.get_ip(enb_serial)

    def _get_handler(
        self,
        client_ip: str,
    ) -> EnodebAcsStateMachine:
        return self._state_machine_by_ip[client_ip]

    def _update_device_mapping(
        self,
        client_ip: str,
        inform: models.Inform,
    ) -> None:
        """
        When receiving an Inform message, we can figure out what device we
        are talking to. We can also see if the IP has changed, and the
        StateMachineManager must track this so that subsequent tr069
        messages can be handled correctly.
        """
        enb_serial = self._parse_msg_for_serial(inform)
        self._associate_serial_to_ip(client_ip, enb_serial)
        handler = self._get_handler(client_ip)
        if handler is None:
            device_name = get_device_name_from_inform(inform)
            handler = self._build_handler(device_name)
            self._state_machine_by_ip[client_ip] = handler

    def _associate_serial_to_ip(
        self,
        client_ip: str,
        enb_serial: str,
    ) -> None:
        """
        If a device/IP combination changes, then the StateMachineManager
        must detect this, and update its mapping of what serial/IP corresponds
        to which handler.
        """
        if enb_serial is None:
            # TR-069 message did not contain an eNodeB serial ID.
            logging.error('Cannot associate null eNB serial to a an IP')
            pass
        elif self._ip_serial_mapping.has_ip(client_ip):
            # Same IP, different eNB connected
            prev_serial = self._ip_serial_mapping.get_serial(client_ip)
            if enb_serial != prev_serial:
                logging.info('eNodeB change on IP <%s>, from %s to %s',
                             client_ip, prev_serial, enb_serial)
                self._ip_serial_mapping.set_ip_and_serial(client_ip, enb_serial)
                self._state_machine_by_ip[client_ip] = None
        elif self._ip_serial_mapping.has_serial(enb_serial):
            # Same eNB, different IP
            prev_ip = self._ip_serial_mapping.get_ip(enb_serial)
            if client_ip != prev_ip:
                logging.info('eNodeB <%s> changed IP from %s to %s',
                             enb_serial, prev_ip, client_ip)
                self._ip_serial_mapping.set_ip_and_serial(client_ip, enb_serial)
                handler = self._state_machine_by_ip[prev_ip]
                self._state_machine_by_ip[client_ip] = handler
                del self._state_machine_by_ip[prev_ip]
        else:
            # TR069 message is coming from a different IP, and a different
            # serial ID. No need to change mapping
            handler = None
            self._ip_serial_mapping.set_ip_and_serial(client_ip, enb_serial)
            self._state_machine_by_ip[client_ip] = handler

    def _parse_msg_for_serial(
        self,
        tr069_message: models.Inform,
    ) -> Optional[str]:
        """ Return the eNodeB serial ID if it's found in the message """
        if not isinstance(tr069_message, models.Inform):
            return
        if not hasattr(tr069_message, 'ParameterList') or \
                not hasattr(tr069_message.ParameterList, 'ParameterValueStruct'):
            return None

        # Parse the parameters
        param_values_by_path = {}
        for param_value in tr069_message.ParameterList.ParameterValueStruct:
            path = param_value.Name
            value = param_value.Value.Data
            param_values_by_path[path] = value

        enb_serial = param_values_by_path['Device.DeviceInfo.SerialNumber']
        return enb_serial

    def _get_client_ip(self, ctx: WsgiMethodContext) -> str:
        return ctx.transport.req_env.get("REMOTE_ADDR", "unknown")

    def _build_handler(
        self,
        device_name: EnodebDeviceName,
    ) -> EnodebAcsStateMachine:
        """
        Create a new state machine based on the device type
        """
        device_handler_class = get_device_handler_from_name(device_name)
        acs_state_machine = device_handler_class(self._service)
        return acs_state_machine


class IpToSerialMapping:
    """ Bidirectional map between <eNodeB IP> and <eNodeB serial ID> """
    def __init__(self) -> None:
        self.ip_by_enb_serial = {}
        self.enb_serial_by_ip = {}

    def del_ip(self, ip: str) -> None:
        if ip not in self.enb_serial_by_ip:
            raise KeyError('Cannot delete missing IP')
        serial = self.enb_serial_by_ip[ip]
        del self.enb_serial_by_ip[ip]
        del self.ip_by_enb_serial[serial]

    def del_serial(self, serial: str) -> None:
        if serial not in self.ip_by_enb_serial:
            raise KeyError('Cannot delete missing eNodeB serial ID')
        ip = self.ip_by_enb_serial[serial]
        del self.ip_by_enb_serial[serial]
        del self.enb_serial_by_ip[ip]

    def set_ip_and_serial(self, ip: str, serial: str) -> None:
        self.ip_by_enb_serial[serial] = ip
        self.enb_serial_by_ip[ip] = serial

    def get_ip(self, serial: str) -> str:
        return self.ip_by_enb_serial[serial]

    def get_serial(self, ip: str) -> str:
        return self.enb_serial_by_ip[ip]

    def has_ip(self, ip: str) -> bool:
        return ip in self.enb_serial_by_ip

    def has_serial(self, serial: str) -> bool:
        return serial in self.ip_by_enb_serial

    def get_serial_list(self) -> List[str]:
        return list(self.ip_by_enb_serial.keys())

    def get_ip_list(self) -> List[str]:
        return list(self.enb_serial_by_ip.keys())