Name:		dante
Version:	0.0.0
Release:	0%{?dist}
Summary:	Single stack reports made simple

License:	GPLv3
URL:	    https://github.com/nethesis/dante	
Source0:	https://github.com/nethesis/dante/archive/master.tar.gz
Source1:    caronte


%description
Single stack reports made simple


%prep
%setup -q -n dante-master


%install
mkdir -p %{buildroot}/usr/share/dante/
mkdir -p %{buildroot}/usr/bin
cp ciacco/ciacco %{buildroot}/%{_bindir}
mv %{SOURCE1}  %{buildroot}/%{_bindir}


%files
%doc README.md
%license LICENSE
/usr/share/dante/
%{_bindir}/ciacco
%{_bindir}/caronte



%changelog

