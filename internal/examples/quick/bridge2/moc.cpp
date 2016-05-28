

#define protected public
#define private public

#include "moc.h"
#include "_cgo_export.h"

#include <QChildEvent>
#include <QEvent>
#include <QMetaMethod>
#include <QMetaObject>
#include <QObject>
#include <QString>
#include <QTime>
#include <QTimer>
#include <QTimerEvent>

class QmlBridge: public QObject
{
Q_OBJECT
public:
	QmlBridge(QObject *parent) : QObject(parent) {};
	void Signal_SendToQml(QString data) { callbackQmlBridge_SendToQml(this, this->objectName().toUtf8().data(), data.toUtf8().data()); };
	void timerEvent(QTimerEvent * event) { callbackQmlBridge_TimerEvent(this, this->objectName().toUtf8().data(), event); };
	void childEvent(QChildEvent * event) { callbackQmlBridge_ChildEvent(this, this->objectName().toUtf8().data(), event); };
	void connectNotify(const QMetaMethod & sign) { callbackQmlBridge_ConnectNotify(this, this->objectName().toUtf8().data(), new QMetaMethod(sign)); };
	void customEvent(QEvent * event) { callbackQmlBridge_CustomEvent(this, this->objectName().toUtf8().data(), event); };
	void deleteLater() { callbackQmlBridge_DeleteLater(this, this->objectName().toUtf8().data()); };
	void disconnectNotify(const QMetaMethod & sign) { callbackQmlBridge_DisconnectNotify(this, this->objectName().toUtf8().data(), new QMetaMethod(sign)); };
	bool event(QEvent * e) { return callbackQmlBridge_Event(this, this->objectName().toUtf8().data(), e) != 0; };
	bool eventFilter(QObject * watched, QEvent * event) { return callbackQmlBridge_EventFilter(this, this->objectName().toUtf8().data(), watched, event) != 0; };
	
signals:
	void sendToQml(QString data);
public slots:
	QString sendToGo(QString data) { return QString(callbackQmlBridge_SendToGo(this, this->objectName().toUtf8().data(), data.toUtf8().data())); };
};

void QmlBridge_ConnectSendToQml(void* ptr)
{
	QObject::connect(static_cast<QmlBridge*>(ptr), static_cast<void (QmlBridge::*)(QString)>(&QmlBridge::sendToQml), static_cast<QmlBridge*>(ptr), static_cast<void (QmlBridge::*)(QString)>(&QmlBridge::Signal_SendToQml));
}

void QmlBridge_DisconnectSendToQml(void* ptr)
{
	QObject::disconnect(static_cast<QmlBridge*>(ptr), static_cast<void (QmlBridge::*)(QString)>(&QmlBridge::sendToQml), static_cast<QmlBridge*>(ptr), static_cast<void (QmlBridge::*)(QString)>(&QmlBridge::Signal_SendToQml));
}

void QmlBridge_SendToQml(void* ptr, char* data)
{
	static_cast<QmlBridge*>(ptr)->sendToQml(QString(data));
}

char* QmlBridge_SendToGo(void* ptr, char* data)
{
	QString returnArg;
	QMetaObject::invokeMethod(static_cast<QmlBridge*>(ptr), "sendToGo", Q_RETURN_ARG(QString, returnArg), Q_ARG(QString, QString(data)));
	return returnArg.toUtf8().data();
}

void* QmlBridge_NewQmlBridge(void* parent)
{
	return new QmlBridge(static_cast<QObject*>(parent));
}

void QmlBridge_DestroyQmlBridge(void* ptr)
{
	static_cast<QmlBridge*>(ptr)->~QmlBridge();
}

void QmlBridge_TimerEvent(void* ptr, void* event)
{
	static_cast<QmlBridge*>(ptr)->timerEvent(static_cast<QTimerEvent*>(event));
}

void QmlBridge_TimerEventDefault(void* ptr, void* event)
{
	static_cast<QmlBridge*>(ptr)->QObject::timerEvent(static_cast<QTimerEvent*>(event));
}

void QmlBridge_ChildEvent(void* ptr, void* event)
{
	static_cast<QmlBridge*>(ptr)->childEvent(static_cast<QChildEvent*>(event));
}

void QmlBridge_ChildEventDefault(void* ptr, void* event)
{
	static_cast<QmlBridge*>(ptr)->QObject::childEvent(static_cast<QChildEvent*>(event));
}

void QmlBridge_ConnectNotify(void* ptr, void* sign)
{
	static_cast<QmlBridge*>(ptr)->connectNotify(*static_cast<QMetaMethod*>(sign));
}

void QmlBridge_ConnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<QmlBridge*>(ptr)->QObject::connectNotify(*static_cast<QMetaMethod*>(sign));
}

void QmlBridge_CustomEvent(void* ptr, void* event)
{
	static_cast<QmlBridge*>(ptr)->customEvent(static_cast<QEvent*>(event));
}

void QmlBridge_CustomEventDefault(void* ptr, void* event)
{
	static_cast<QmlBridge*>(ptr)->QObject::customEvent(static_cast<QEvent*>(event));
}

void QmlBridge_DeleteLater(void* ptr)
{
	QMetaObject::invokeMethod(static_cast<QmlBridge*>(ptr), "deleteLater");
}

void QmlBridge_DeleteLaterDefault(void* ptr)
{
	static_cast<QmlBridge*>(ptr)->QObject::deleteLater();
}

void QmlBridge_DisconnectNotify(void* ptr, void* sign)
{
	static_cast<QmlBridge*>(ptr)->disconnectNotify(*static_cast<QMetaMethod*>(sign));
}

void QmlBridge_DisconnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<QmlBridge*>(ptr)->QObject::disconnectNotify(*static_cast<QMetaMethod*>(sign));
}

int QmlBridge_Event(void* ptr, void* e)
{
	return static_cast<QmlBridge*>(ptr)->event(static_cast<QEvent*>(e));
}

int QmlBridge_EventDefault(void* ptr, void* e)
{
	return static_cast<QmlBridge*>(ptr)->QObject::event(static_cast<QEvent*>(e));
}

int QmlBridge_EventFilter(void* ptr, void* watched, void* event)
{
	return static_cast<QmlBridge*>(ptr)->eventFilter(static_cast<QObject*>(watched), static_cast<QEvent*>(event));
}

int QmlBridge_EventFilterDefault(void* ptr, void* watched, void* event)
{
	return static_cast<QmlBridge*>(ptr)->QObject::eventFilter(static_cast<QObject*>(watched), static_cast<QEvent*>(event));
}





#include "moc_moc.h"
